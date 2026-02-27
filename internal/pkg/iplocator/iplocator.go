package iplocator

import (
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/pkg/respath"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lionsoul2014/ip2region/binding/golang/service"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/mlogclub/simple/common/strs"
)

const (
	IPv4DataURL = "https://github.com/lionsoul2014/ip2region/raw/refs/heads/master/data/ip2region_v4.xdb"
	IPv6DataURL = "https://github.com/lionsoul2014/ip2region/raw/refs/heads/master/data/ip2region_v6.xdb"
)

var (
	once          sync.Once
	ip2regionAtom atomic.Pointer[service.Ip2Region]
)

func InitIpLocator() {
	once.Do(func() {
		go initIpLocatorAsync()
	})
}

func Search(ip string) string {
	svc := ip2regionAtom.Load()
	if svc == nil || strs.IsBlank(ip) {
		return ""
	}
	region, err := svc.SearchByStr(ip)
	if err != nil {
		return ""
	}
	return region
}

func IpLocation(ip string) string {
	region := Search(ip) // eg. 中国|0|湖北省|武汉市|电信
	return ipLocationFromRegion(region)
}

func ipLocationFromRegion(region string) string {
	if strs.IsBlank(region) {
		return ""
	}
	ss := strings.Split(region, "|")
	if len(ss) < 2 {
		return ""
	}
	var (
		nation   = normalizeRegionPart(ss[0])
		province = pickProvince(ss)
	)
	if strs.IsNotBlank(province) {
		return province
	}
	if strs.IsNotBlank(nation) {
		return nation
	}
	return ""
}

func verifyXdb(path string, version string) string {
	if strs.IsBlank(path) {
		return ""
	}
	if err := xdb.VerifyFromFile(path); err != nil {
		slog.Error("verify ip2region xdb failed, disable this version", slog.Any("version", version), slog.Any("path", path), slog.Any("err", err))
		return ""
	}
	return path
}

func pickProvince(ss []string) string {
	if len(ss) >= 2 {
		if v := normalizeRegionPart(ss[1]); strs.IsNotBlank(v) {
			return v
		}
	}
	if len(ss) >= 3 {
		if v := normalizeRegionPart(ss[2]); strs.IsNotBlank(v) {
			return v
		}
	}
	return ""
}

func normalizeRegionPart(part string) string {
	part = strings.TrimSpace(part)
	if part == "" || part == "0" {
		return ""
	}
	return part
}

func initIpLocatorAsync() {
	// start := dates.NowTimestamp()
	v4Path, v6Path := resolveDataPaths()
	v4Path = verifyXdb(v4Path, "v4")
	v6Path = verifyXdb(v6Path, "v6")

	if strs.IsBlank(v4Path) && strs.IsBlank(v6Path) {
		slog.Error("ip2region disabled, both v4 and v6 data are unavailable")
		return
	}

	svc, err := service.NewIp2RegionWithPath(v4Path, v6Path)
	if err != nil {
		slog.Error("failed to create ip2region service", slog.Any("v4Path", v4Path), slog.Any("v6Path", v6Path), slog.Any("err", err))
		return
	}
	ip2regionAtom.Store(svc)
	// slog.Info("load ip2region service success", slog.Any("v4Path", v4Path), slog.Any("v6Path", v6Path), slog.Any("elapsed", dates.NowTimestamp()-start))
}

func resolveDataPaths() (string, string) {
	v4Path := strings.TrimSpace(config.Instance.IPLocator.IPv4DataPath)
	v6Path := strings.TrimSpace(config.Instance.IPLocator.IPv6DataPath)

	// When both paths are absent, bootstrap default data into res/uploads.
	if strs.IsBlank(v4Path) && strs.IsBlank(v6Path) {
		v4Path = ensureDataFile(respath.UploadsPath("ip2region_v4.xdb"), IPv4DataURL, "v4")
		v6Path = ensureDataFile(respath.UploadsPath("ip2region_v6.xdb"), IPv6DataURL, "v6")
	}
	return v4Path, v6Path
}

func ensureDataFile(localPath string, dataURL string, version string) string {
	if isValidXdb(localPath) {
		// slog.Info("ip2region data file ready", slog.Any("version", version), slog.Any("path", localPath))
		return localPath
	}
	if err := downloadFile(localPath, dataURL, version); err != nil {
		slog.Warn("download ip2region data failed", slog.Any("version", version), slog.Any("url", dataURL), slog.Any("path", localPath), slog.Any("err", err))
		return ""
	}
	if !isValidXdb(localPath) {
		slog.Warn("downloaded ip2region data invalid", slog.Any("version", version), slog.Any("path", localPath))
		return ""
	}
	return localPath
}

func isValidXdb(path string) bool {
	if strs.IsBlank(path) {
		return false
	}
	return xdb.VerifyFromFile(path) == nil
}

func downloadFile(localPath string, url string, version string) error {
	if strs.IsBlank(localPath) || strs.IsBlank(url) {
		return fmt.Errorf("invalid download args")
	}
	if err := os.MkdirAll(filepath.Dir(localPath), 0o755); err != nil {
		return err
	}
	startAt := time.Now()
	// slog.Info("start downloading ip2region data", slog.Any("version", version), slog.Any("url", url), slog.Any("path", localPath))

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %s", resp.Status)
	}
	// slog.Info("ip2region data response received", slog.Any("version", version), slog.Any("status", resp.Status), slog.Any("contentLength", resp.ContentLength))

	tmpPath := fmt.Sprintf("%s.tmp.%d", localPath, time.Now().UnixNano())
	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	progress := newDownloadProgressReader(resp.Body, version, resp.ContentLength)
	written, copyErr := io.Copy(tmpFile, progress)
	progress.logProgress(true)
	closeErr := tmpFile.Close()
	if copyErr != nil {
		_ = os.Remove(tmpPath)
		return copyErr
	}
	if closeErr != nil {
		_ = os.Remove(tmpPath)
		return closeErr
	}

	if err = os.Rename(tmpPath, localPath); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	slog.Info("ip2region data download completed", slog.Any("version", version), slog.Any("path", localPath), slog.Any("bytes", written), slog.Any("elapsedMs", time.Since(startAt).Milliseconds()))
	return nil
}

type downloadProgressReader struct {
	reader          io.Reader
	version         string
	totalBytes      int64
	downloadedBytes int64
	lastLogBytes    int64
	lastLogAt       time.Time
	startAt         time.Time
	inlineProgress  bool
}

func newDownloadProgressReader(reader io.Reader, version string, totalBytes int64) *downloadProgressReader {
	now := time.Now()
	return &downloadProgressReader{
		reader:         reader,
		version:        version,
		totalBytes:     totalBytes,
		lastLogAt:      now,
		startAt:        now,
		inlineProgress: isTerminal(os.Stdout),
	}
}

func (r *downloadProgressReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	if n > 0 {
		r.downloadedBytes += int64(n)
		r.logProgress(false)
	}
	return n, err
}

func (r *downloadProgressReader) logProgress(force bool) {
	const (
		minLogBytes    = int64(2 * 1024 * 1024)
		minLogInterval = 2 * time.Second
	)
	now := time.Now()
	if !force {
		if r.downloadedBytes-r.lastLogBytes < minLogBytes && now.Sub(r.lastLogAt) < minLogInterval {
			return
		}
	}

	elapsed := now.Sub(r.startAt)
	if elapsed <= 0 {
		elapsed = time.Millisecond
	}
	speed := float64(r.downloadedBytes) / elapsed.Seconds()
	if r.inlineProgress {
		if r.totalBytes > 0 {
			percent := float64(r.downloadedBytes) * 100 / float64(r.totalBytes)
			fmt.Printf("\rip2region %s downloading: %d/%d bytes (%.2f%%) speed=%d B/s", r.version, r.downloadedBytes, r.totalBytes, percent, int64(speed))
		} else {
			fmt.Printf("\rip2region %s downloading: %d bytes speed=%d B/s", r.version, r.downloadedBytes, int64(speed))
		}
		if force {
			fmt.Println()
		}
	}
	r.lastLogBytes = r.downloadedBytes
	r.lastLogAt = now
}

func isTerminal(file *os.File) bool {
	if file == nil {
		return false
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}
