const { request } = require('../../utils/request');
const { fen2yuan, toast } = require('../../utils/util');

Page({
  data: { month: '', summary: null, loading: true },

  onShow() {
    const d = new Date();
    const month = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`;
    this.setData({ month }, () => this.load());
  },

  shift(e) {
    const delta = Number(e.currentTarget.dataset.d);
    const [y, m] = this.data.month.split('-').map(Number);
    const d = new Date(y, m - 1 + delta, 1);
    const month = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`;
    this.setData({ month }, () => this.load());
  },

  load() {
    this.setData({ loading: true });
    return request(`/income/summary?month=${this.data.month}`)
      .then(s => this.setData({
        summary: {
          totalYuan: fen2yuan(s.total),
          orders: s.orders,
          hours: (s.hours || 0).toFixed(1),
          byStudent: (s.by_student || []).map(x => ({ ...x, totalYuan: fen2yuan(x.total) })),
        },
        loading: false,
      }))
      .catch(err => {
        this.setData({ loading: false });
        toast(err.message);
      });
  },
});
