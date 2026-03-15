package migrations

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/dto"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/sqls"
)

func migrate_site_content_config() error {
	return sqls.WithTransaction(func(ctx *sqls.TxContext) error {
		tx := ctx.Tx
		now := dates.NowTimestamp()

		if repositories.SysConfigRepository.GetByKey(tx, constants.SysConfigAboutPageConfig) == nil {
			value, err := jsons.ToStr(defaultAboutPageConfig())
			if err != nil {
				return err
			}
			name, desc := aboutPageConfigMetaByLanguage()
			if err := repositories.SysConfigRepository.Create(tx, &models.SysConfig{
				Key:         constants.SysConfigAboutPageConfig,
				Value:       value,
				Name:        name,
				Description: desc,
				CreateTime:  now,
				UpdateTime:  now,
			}); err != nil {
				return err
			}
		}

		if repositories.SysConfigRepository.GetByKey(tx, constants.SysConfigFooterLinks) == nil {
			value, err := jsons.ToStr(defaultFooterLinks())
			if err != nil {
				return err
			}
			name, desc := footerLinksMetaByLanguage()
			if err := repositories.SysConfigRepository.Create(tx, &models.SysConfig{
				Key:         constants.SysConfigFooterLinks,
				Value:       value,
				Name:        name,
				Description: desc,
				CreateTime:  now,
				UpdateTime:  now,
			}); err != nil {
				return err
			}
		}
		return nil
	})
}

func defaultAboutPageConfig() dto.AboutPageConfig {
	return dto.AboutPageConfig{
		Content: dto.LocalizedText{
			"en-US": "# About Us\n\nBBS-GO is a modern, high-performance open source community forum system. Our design philosophy is lightweight, efficient, easy to extend and deploy, aiming to provide developers and community managers with a powerful online community solution.\n\n## Core Features\n\n- **High Performance:** Ensures smooth user experience even under high load.\n- **Highly Flexible:** Supports rich custom configurations and plugin extensions to easily meet business needs in different scenarios.\n- **Easy to Use:** A simple and powerful admin panel makes community management easy and efficient.\n- **Stable & Reliable:** Ensures system stability and good scalability.\n- **Responsive Design:** Perfectly adapts to desktop and mobile devices, providing a consistent experience.\n\n## Join Us\n\nBBS-GO is a vibrant open source project. We welcome all forms of contribution, including code, feature suggestions, documentation, and bug reports.\n\n- **GitHub:** [mlogclub/bbs-go](https://github.com/mlogclub/bbs-go)\n- **Gitee:** [mlogclub/bbs-go](https://gitee.com/mlogclub/bbs-go)\n- **Contributors:** [Contributors List](https://github.com/mlogclub/bbs-go/graphs/contributors)\n\n## Open Source License\n\nBBS-GO follows the [GNU General Public License v3.0](https://github.com/mlogclub/bbs-go/blob/master/LICENSE) open source license.\n\n## Contact & Support\n\nIf you encounter any problems or have business cooperation needs, please contact us through the following ways.\n\n- **Community:** [BBS-GO Community](https://bbs.bbs-go.com)\n- **Docs:** [bbs-go.com](https://bbs-go.com)\n",
			"zh-CN": "# 关于我们\n\nBBS-GO 是一款现代化的、高性能的开源社区论坛系统。我们的设计哲学是轻量、高效、易于扩展和部署，旨在为开发者和社区管理者提供一个强大的在线社区解决方案。\n\n## 核心特性\n\n- **高性能：** 确保在高负载下也能提供流畅的用户体验。\n- **高度灵活：** 支持丰富的自定义配置和插件扩展，轻松满足不同场景的业务需求。\n- **简单易用：** 拥有设计简洁、功能强大的管理后台，让社区管理变得轻松高效。\n- **稳定可靠：** 确保系统稳定性和良好的可扩展性。\n- **响应式设计：** 完美适配桌面和移动设备，为用户提供一致的访问体验。\n\n## 加入我们\n\nBBS-GO 是一个充满活力的开源项目，我们欢迎任何形式的贡献。无论是代码实现、功能建议、文档完善还是 Bug 反馈，都是对我们社区的宝贵支持。\n\n- **GitHub：** [mlogclub/bbs-go](https://github.com/mlogclub/bbs-go)\n- **Gitee：** [mlogclub/bbs-go](https://gitee.com/mlogclub/bbs-go)\n- **项目贡献者：** [贡献者列表](https://github.com/mlogclub/bbs-go/graphs/contributors)\n\n## 开源协议\n\nBBS-GO 遵循 [GNU General Public License v3.0](https://github.com/mlogclub/bbs-go/blob/master/LICENSE) 开源协议。\n\n## 联系与支持\n\n如果您在使用中遇到任何问题，或有任何商业合作需求，欢迎通过以下方式联系我们。\n\n- **交流社区：** [BBS-GO 社区](https://bbs.bbs-go.com)\n- **官方文档：** [bbs-go.com](https://bbs-go.com)\n",
		},
	}
}

func defaultFooterLinks() []dto.FooterLink {
	return []dto.FooterLink{
		{
			Text: dto.LocalizedText{
				"en-US": "About",
				"zh-CN": "关于",
			},
			Url:             "/about",
			OpenInNewWindow: false,
			Visible:         true,
		},
		{
			Text: dto.LocalizedText{
				"en-US": "Links",
				"zh-CN": "友情链接",
			},
			Url:             "/links",
			OpenInNewWindow: false,
			Visible:         true,
		},
		{
			Text: dto.LocalizedText{
				"en-US": "ICP Filing",
				"zh-CN": "ICP备案号",
			},
			Url:             "",
			OpenInNewWindow: true,
			Visible:         false,
		},
		{
			Text: dto.LocalizedText{
				"en-US": "Public Security Filing",
				"zh-CN": "公安网备信息",
			},
			Url:             "",
			OpenInNewWindow: true,
			Visible:         false,
		},
	}
}

func aboutPageConfigMetaByLanguage() (name, description string) {
	if config.Instance.Language == config.LanguageEnUS {
		return "About Page Config", "Configurable about page content with i18n text"
	}
	return "关于页配置", "可配置的关于页内容，支持中英文文本"
}

func footerLinksMetaByLanguage() (name, description string) {
	if config.Instance.Language == config.LanguageEnUS {
		return "Footer Links", "Footer links and filing information configuration"
	}
	return "底部链接", "底部链接与备案信息等配置"
}
