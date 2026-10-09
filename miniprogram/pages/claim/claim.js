const { request } = require('../../utils/request');
const { toast } = require('../../utils/util');

Page({
  data: { code: '' },

  onInput(e) {
    this.setData({ code: e.detail.value });
  },

  onClaim() {
    const code = this.data.code.trim();
    if (!/^\d{6}$/.test(code)) return toast('认领码是 6 位数字');
    request('/coach/claim', { method: 'POST', data: { claim_code: code } })
      .then(({ token }) => {
        wx.setStorageSync('token', token);
        wx.showModal({
          title: '认领成功',
          content: '档案已绑定到你的微信，去完善主页吧',
          showCancel: false,
          success: () => wx.switchTab({ url: '/pages/me/me' }),
        });
      })
      .catch(err => toast(err.message));
  },
});
