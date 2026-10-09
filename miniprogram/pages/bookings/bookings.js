const { request } = require('../../utils/request');
const { STATUS_TEXT, fen2yuan, toast } = require('../../utils/util');

const TABS = [
  { k: '', t: '全部' },
  { k: 'pending', t: '待确认' },
  { k: 'confirmed', t: '已确认' },
  { k: 'completed', t: '已完成' },
  { k: 'cancelled', t: '已取消' },
];

Page({
  data: { tabs: TABS, active: '', list: [], total: 0, page: 1, loading: false, noMore: false },

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

  switchTab(e) {
    this.setData({ active: e.currentTarget.dataset.k, page: 1, noMore: false }, () => this.load(true));
  },

  load(reset) {
    const { active, page } = this.data;
    this.setData({ loading: true });
    return request(`/bookings?status=${active}&page=${page}&page_size=20`)
      .then(({ list, total }) => {
        const mapped = (list || []).map(b => ({
          ...b,
          statusText: STATUS_TEXT[b.status] || b.status,
          priceYuan: fen2yuan(b.price),
          dateLabel: `${b.play_date} ${b.start_time.slice(0, 5)}–${b.end_time.slice(0, 5)}`,
        }));
        this.setData({
          list: reset ? mapped : this.data.list.concat(mapped),
          total,
          loading: false,
          noMore: (reset ? mapped.length : this.data.list.length + mapped.length) >= total,
        });
      })
      .catch(err => {
        this.setData({ loading: false });
        toast(err.message);
      });
  },

  reload(id, fn) {
    fn(id)
      .then(() => { this.setData({ page: 1, noMore: false }); return this.load(true); })
      .catch(err => toast(err.message));
  },

  onConfirm(e) {
    this.reload(e.currentTarget.dataset.id, id => request(`/bookings/${id}/confirm`, { method: 'PATCH' }));
  },

  onReject(e) {
    const id = e.currentTarget.dataset.id;
    wx.showModal({
      title: '拒绝这单？', editable: true, placeholderText: '理由（选填）',
      success: res => {
        if (!res.confirm) return;
        this.reload(id, i => request(`/bookings/${i}/reject`, { method: 'PATCH', data: { reason: res.content || '' } }));
      },
    });
  },

  onComplete(e) {
    const id = e.currentTarget.dataset.id;
    wx.showModal({
      title: '标记完成？', content: '完成后学员可以评价',
      success: res => {
        if (res.confirm) this.reload(id, i => request(`/bookings/${i}/complete`, { method: 'PATCH' }));
      },
    });
  },

  onCancel(e) {
    const id = e.currentTarget.dataset.id;
    wx.showModal({
      title: '取消这单？', editable: true, placeholderText: '取消原因（选填）',
      success: res => {
        if (!res.confirm) return;
        this.reload(id, i => request(`/bookings/${i}/cancel`, { method: 'PATCH', data: { reason: res.content || '', no_show_by: 'coach' } }));
      },
    });
  },

  onMarkPaid(e) {
    const id = e.currentTarget.dataset.id;
    wx.showModal({
      title: '确认收款？', content: '确认后计入收入',
      success: res => {
        if (res.confirm) this.reload(id, i => request(`/bookings/${i}/mark-paid`, { method: 'PATCH' }));
      },
    });
  },
});
