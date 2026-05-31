package req

import (
	"bbs-go/internal/pkg/idcodec"
	"testing"
)

func TestUserReportReqDecodedDataIdAcceptsEncodedId(t *testing.T) {
	idcodec.Init(1)
	encoded := idcodec.Encode(12345)

	req := UserReportReq{DataId: encoded}

	if got := req.DecodedDataId(); got != 12345 {
		t.Fatalf("DecodedDataId() = %d, want 12345", got)
	}
}

func TestUserReportReqDecodedDataIdAcceptsNumericId(t *testing.T) {
	idcodec.Init(1)

	req := UserReportReq{DataId: "12345"}

	if got := req.DecodedDataId(); got != 12345 {
		t.Fatalf("DecodedDataId() = %d, want 12345", got)
	}
}
