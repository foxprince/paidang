const { request } = require('../../utils/request');
const { fen2yuan, toast } = require('../../utils/util');

const STATUS_TEXT = { active: '正常', suspended: '已下架', banned: '已封号' };

Page({
  data: { id: null, coach: null, bookings: 0, reviews: 0, statusText: '' },

  onLoad(options) {
    this.setData({ id: options.id });
    this.load();
  },

  load() {
    return request(`/admin/coaches/${this.data.id}`)
      .then(({ coach, bookings, reviews }) => this.setData({
        coach: { ...coach, priceYuan: fen2yuan(coach.price_per_hour) },
        bookings, reviews,
        statusText: STATUS_TEXT[coach.status] || coach.status,
      }))
      .catch(err => toast(err.message));
  },

  askNote(title, cb) {
    wx.showModal({
      title, editable: true, placeholderText: '备注（会记入审计日志）',
      success: res => { if (res.confirm) cb(res.content || ''); },
    });
  },

  onCopyCode() {
    wx.setClipboardData({ data: this.data.coach.claim_code || '' });
  },

  onVerify(e) {    const verified = e.currentTarget.dataset.v === '1';
    this.askNote(verified ? '通过审核？' : '驳回并下架？', note => {
      request(`/admin/coaches/${this.data.id}/verify`, { method: 'PATCH', data: { verified, note } })
        .then(() => { toast(verified ? '已通过' : '已下架'); return this.load(); })
        .catch(err => toast(err.message));
    });
  },

  onStatus(e) {
    const status = e.currentTarget.dataset.s;
    const title = { suspended: '下架这个陪练？', banned: '封号？请慎重', active: '恢复正常？' }[status];
    wx.showModal({
      title, content: status === 'banned' ? '封号后不可登录，不可逆操作请三思' : '',
      success: res => {
        if (!res.confirm) return;
        this.askNote('备注', note => {
          request(`/admin/coaches/${this.data.id}/status`, { method: 'PATCH', data: { status, note } })
            .then(() => { toast('已执行'); return this.load(); })
            .catch(err => toast(err.message));
        });
      },
    });
  },
});
