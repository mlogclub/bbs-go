export default {
  settings: {
    common: {
      title: 'General',
      siteTitle: 'Site Title',
      siteLogo: 'Site Logo',
      siteDescription: 'Site Description',
      siteKeywords: 'Site Keywords',
      siteNotification: 'Site Notification',
      recommendTags: 'Recommended Tags',
      defaultNodeId: 'Default Node',
      modules: 'Function Modules',
      tweet: 'Tweet',
      topic: 'Topic',
      article: 'Article',
      urlRedirect: 'External Link Redirect Page',
      urlRedirectTooltip:
        'Manual confirmation required before redirecting to external links',
      enableHideContent: 'Enable Comment-Visible Content',
      enableHideContentTooltip:
        'Support setting content visible after commenting when posting',
      submit: 'Submit',
      placeholder: {
        siteTitle: 'Site Title',
        siteDescription: 'Site Description',
        siteKeywords: 'Site Keywords',
        siteNotification: 'Site Notification (HTML supported)',
        recommendTags: 'Recommended Tags',
        defaultNodeId: 'Default node for posting',
      },
      message: {
        submitSuccess: 'Submit successful',
      },
    },
    nav: {
      title: 'Navigation',
      tableTitle: 'Title',
      tableUrl: 'URL',
      submit: 'Submit',
      message: {
        validation:
          'Please check your navigation configuration, title and URL cannot be empty',
        submitSuccess: 'Submit successful',
      },
    },
    score: {
      title: 'Score',
      postTopicScore: 'Post Topic Score',
      postCommentScore: 'Post Comment Score',
      checkInScore: 'Check-in Score',
      submit: 'Submit',
      placeholder: {
        postTopicScore: 'Score earned for posting topic',
        postCommentScore: 'Score earned for posting comment',
        checkInScore: 'Score earned for check-in',
      },
      message: {
        submitSuccess: 'Submit successful',
      },
    },
    spam: {
      title: 'Anti-Spam',
      topicCaptcha: 'Topic Captcha',
      topicCaptchaTooltip:
        'Whether to enable captcha verification when posting topics',
      createTopicEmailVerified: 'Email Verified for Topic Creation',
      createTopicEmailVerifiedTooltip:
        'Email verification required before posting topics',
      createArticleEmailVerified: 'Email Verified for Article Creation',
      createArticleEmailVerifiedTooltip:
        'Email verification required before posting articles',
      createCommentEmailVerified: 'Email Verified for Comments',
      createCommentEmailVerifiedTooltip:
        'Email verification required before posting comments',
      articlePending: 'Article Review',
      articlePendingTooltip:
        'Whether to enable review after publishing articles',
      userObserveSeconds: 'User Observation Period (seconds)',
      userObserveSecondsTooltip:
        'During the observation period, users cannot post topics, tweets, etc. Set to 0 for no observation period.',
      emailWhitelist: 'Email Whitelist',
      submit: 'Submit',
      placeholder: {
        emailWhitelist: 'Email Whitelist',
      },
      message: {
        submitSuccess: 'Submit successful',
      },
    },
    upload: {
      title: 'Upload',
      uploadConfig: 'Upload Configuration',
      enableUploadMethod: 'Upload Method',
      aliyunOss: 'Alibaba Cloud OSS',
      tencentCos: 'Tencent Cloud COS',
      host: 'Domain',
      bucket: 'Bucket',
      endpoint: 'Endpoint',
      accessKeyId: 'AccessKey ID',
      accessKeySecret: 'AccessKey Secret',
      region: 'Region',
      secretId: 'SecretId',
      secretKey: 'SecretKey',
      imageStyleConfig: 'Image Style Configuration',
      styleSplitter: 'Style Splitter',
      styleAvatar: 'Avatar Style',
      stylePreview: 'Preview Style',
      styleSmall: 'Small Image Style',
      styleDetail: 'Detail Style',
      submit: 'Submit',
      placeholder: {
        host: 'Enter OSS domain, e.g.: https://xxx.oss-cn-beijing.aliyuncs.com/',
        bucket: 'Enter Bucket name',
        endpoint: 'Enter Endpoint, e.g.: oss-cn-beijing.aliyuncs.com',
        accessKeyId: 'Enter AccessKey ID',
        accessKeySecret: 'Enter AccessKey Secret',
        tencentBucket: 'Enter Bucket name, format: bucket-appid',
        region: 'Enter region, e.g.: ap-beijing',
        secretId: 'Enter SecretId',
        secretKey: 'Enter SecretKey',
        styleSplitter: 'Style splitter, e.g.: !',
        styleAvatar: 'Avatar style, e.g.: 100x100',
        stylePreview: 'Preview style, e.g.: 400x400',
        styleSmall: 'Small image style, e.g.: 200x200',
        styleDetail: 'Detail style, e.g.: 800x800',
      },
      message: {
        submitSuccess: 'Submit successful',
      },
    },
  },
};
