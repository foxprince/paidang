const { request } = require('../../utils/request');
const { STATUS_TEXT, toast } = require('../../utils/util');

Page({
  data: { list: [], loading: true },

  onShow() {
    this.load();
  },

  onPullDownRefresh() {
    this.load().finally(() => wx.stopPullDownRefresh());
  },

  load() {
    this.setData({ loading: true });
    return request('/coach/today')
      .then(list => this.setData({
        list: (list || []).map(b => ({
          ...b,
          statusText: STATUS_TEXT[b.status] || b.status,
          timeLabel: `${b.start_time.slice(0, 5)}–${b.end_time.slice(0, 5)}`,
        })),
        loading: false,
      }))
      .catch(err => {
        this.setData({ loading: false });
        toast(err.message);
      });
  },

  reload(id, fn) {
    fn(id).then(() => this.load()).catch(err => toast(err.message));
  },

  onConfirm(e) {
    this.reload(e.currentTarget.dataset.id, id => request(`/bookings/${id}/confirm`, { method: 'PATCH' }));
  },

  onComplete(e) {
    const id = e.currentTarget.dataset.id;
    wx.showModal({
      title: '标记完成？',
      content: '完成后学员可以评价',
      success: res => {
        if (res.confirm) this.reload(id, i => request(`/bookings/${i}/complete`, { method: 'PATCH' }));
      },
    });
  },

  onReject(e) {
    const id = e.currentTarget.dataset.id;
    wx.showModal({
      title: '拒绝这单？',
      editable: true,
      placeholderText: '理由（选填）',
      success: res => {
        if (!res.confirm) return;
        this.reload(id, i => request(`/bookings/${i}/reject`, { method: 'PATCH', data: { reason: res.content || '' } }));
      },
    });
  },
});
