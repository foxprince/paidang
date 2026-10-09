const { request } = require('../../utils/request');
const { fen2yuan, toast } = require('../../utils/util');

Page({
  data: { id: null, detail: null },

  onLoad(options) {
    this.setData({ id: options.id });
    this.load();
  },

  load() {
    return request(`/admin/complaints/${this.data.id}`)
      .then(d => {
        let evidence = [];
        try { evidence = JSON.parse(d.complaint.evidence || '[]'); } catch (e) {}
        this.setData({
          detail: {
            ...d,
            evidence,
            priceYuan: fen2yuan(d.booking.price),
            timeLabel: `${d.booking.play_date} ${d.booking.start_time.slice(0, 5)}–${d.booking.end_time.slice(0, 5)}`,
          },
        });
      })
      .catch(err => toast(err.message));
  },

  judge(e) {
    const result = e.currentTarget.dataset.r; // upheld / rejected
    const title = result === 'upheld' ? '裁决成立？' : '驳回投诉？';
    const tip = result === 'upheld' ? '成立将给陪练记一次履约污点' : '';
    wx.showModal({
      title, content: tip, editable: true, placeholderText: '裁决备注（记入审计日志）',
      success: res => {
        if (!res.confirm) return;
        request(`/admin/complaints/${this.data.id}`, { method: 'PATCH', data: { result, note: res.content || '' } })
          .then(() => { toast('已裁决'); wx.navigateBack(); })
          .catch(err => toast(err.message));
      },
    });
  },
});
