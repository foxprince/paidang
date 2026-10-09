const { request } = require('../../utils/request');
const { toast } = require('../../utils/util');

const TABS = [
  { k: 'false', t: '待审核' },
  { k: '', t: '全部' },
];

Page({
  data: { tabs: TABS, active: 'false', list: [], total: 0 },

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
    return request(`/admin/coaches?verified=${this.data.active}&page=1&page_size=50`)
      .then(({ list, total }) => this.setData({ list: list || [], total }))
      .catch(err => toast(err.message));
  },

  goDetail(e) {
    wx.navigateTo({ url: `/pages/admin-coach-detail/admin-coach-detail?id=${e.currentTarget.dataset.id}` });
  },

  goNew() {
    wx.navigateTo({ url: '/pages/admin-coach-new/admin-coach-new' });
  },
});
