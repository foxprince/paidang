const { request } = require('./utils/request');

App({
  onLaunch() {
    this.silentLogin();
  },

  // 静默登录：拿 code 换 token，存起来
  silentLogin() {
    wx.login({
      success: ({ code }) => {
        request('/auth/wechat-login', {
          method: 'POST',
          data: { code },
          auth: false,
        })
          .then(({ token }) => wx.setStorageSync('token', token))
          .catch(() => {});
      },
    });
  },
});
