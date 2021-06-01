class CommonHelper {
  isMobile(ua) {
    return /mobile|android|webos|iphone|blackberry|micromessenger/i.test(ua)
  }
}

export default new CommonHelper()
