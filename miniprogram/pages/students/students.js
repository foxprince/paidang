const { request } = require('../../utils/request');
const { toast } = require('../../utils/util');

Page({
  data: { list: [], total: 0, page: 1, loading: false, noMore: false },

  onShow() {
    this.setData({ page: 1, noMore: false }, () => this.load(true));
  },

  onPullDownRefresh() {
    this.setData({ page: 1, noMore: false });
    this.load(true).finally(() => wx.stopPullDownRefresh());
  },

  onReachBottom() {
    if (!this.data.noMore && !this.data.loading) {
      this.setData({ page: this.data.page + 1 }, () => this.load(false));
    }
  },

  load(reset) {
    const { page } = this.data;
    this.setData({ loading: true });
    return request(`/students?page=${page}&page_size=20`)
      .then(({ list, total }) => {
        const next = reset ? list : this.data.list.concat(list);
        this.setData({ list: next || [], total, loading: false, noMore: next.length >= total });
      })
      .catch(err => {
        this.setData({ loading: false });
        toast(err.message);
      });
  },

  goDetail(e) {
    wx.navigateTo({ url: `/pages/student-detail/student-detail?id=${e.currentTarget.dataset.id}` });
  },

  goNew() {
    wx.navigateTo({ url: '/pages/student-detail/student-detail?id=new' });
  },
});
