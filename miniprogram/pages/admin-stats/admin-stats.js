const { request } = require('../../utils/request');
const { fen2yuan, toast } = require('../../utils/util');

Page({
  data: { stats: null },

  onShow() {
    request('/admin/stats')
      .then(s => this.setData({
        stats: {
          ...s,
          gmvTotalYuan: fen2yuan(s.gmv_total),
          gmvMonthYuan: fen2yuan(s.gmv_month),
          byStatus: Object.entries(s.bookings_by_status || {}).map(([k, v]) => ({
            k, v, t: { pending: '待确认', confirmed: '已确认', completed: '已完成', cancelled: '已取消' }[k] || k,
          })),
        },
      }))
      .catch(err => toast(err.message));
  },
});
