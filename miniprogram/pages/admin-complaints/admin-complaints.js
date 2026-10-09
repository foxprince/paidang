const { request } = require('../../utils/request');
const { toast } = require('../../utils/util');

const STATUS_TEXT = { open: '待处理', upheld: '成立', rejected: '不成立' };

Page({
  data: { tabs: [{ k: 'open', t: '待处理' }, { k: '', t: '全部' }], active: 'open', list: [] },

  onShow() {
    this.load();
  },

  onPullDownRefresh() {
    this.load().finally(() => wx.stopPullDownRefresh());
  },

  switchTab(e) {
    this.setData({ active: e.currentTarget.dataset.k }, () => this.load());
  },

  load() {
    return request(`/admin/complaints?status=${this.data.active}&page=1&page_size=50`)
      .then(({ list }) => this.setData({
        list: (list || []).map(x => ({ ...x, statusText: STATUS_TEXT[x.status] || x.status })),
      }))
      .catch(err => toast(err.message));
  },

  goDetail(e) {
    wx.navigateTo({ url: `/pages/admin-complaint-detail/admin-complaint-detail?id=${e.currentTarget.dataset.id}` });
  },
});
