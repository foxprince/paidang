const { request } = require('../../utils/request');
const { toast } = require('../../utils/util');

Page({
  data: { stats: null },

  onShow() {
    request('/admin/stats')
      .then(stats => this.setData({ stats }))
      .catch(() => this.setData({ stats: null }));
  },

  go(e) {
    wx.navigateTo({ url: e.currentTarget.dataset.url });
  },
});
