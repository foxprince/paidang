const { request } = require('../../utils/request');
const { TPL_IDS } = require('../../utils/config');
const { nextDays, uuid, toast } = require('../../utils/util');

Page({
  data: {
    id: null,
    dates: [], dateIdx: 0,
    slots: [], slotIdx: -1,
    name: '', phone: '', code: '',
    counting: 0, submitting: false, done: false,
  },

  onLoad(options) {
    this.setData({ id: options.id, dates: nextDays(14) }, () => this.loadSlots());
  },

  pickDate(e) {
    this.setData({ dateIdx: e.currentTarget.dataset.i, slotIdx: -1 }, () => this.loadSlots());
  },

  loadSlots() {
    const date = this.data.dates[this.data.dateIdx].date;
    request(`/public/coach/${this.data.id}/slots?date=${date}`, { auth: false })
      .then(slots => this.setData({
        slots: (slots || []).map(s => ({ ...s, label: `${s.start.slice(0, 5)}–${s.end.slice(0, 5)}` })),
      }))
      .catch(err => toast(err.message));
  },

  pickSlot(e) {
    this.setData({ slotIdx: e.currentTarget.dataset.i });
  },

  onInput(e) {
    this.setData({ [e.currentTarget.dataset.k]: e.detail.value });
  },

  sendCode() {
    const { phone, counting } = this.data;
    if (counting > 0) return;
    if (!/^1\d{10}$/.test(phone)) return toast('手机号格式不对');
    request('/public/sms-code', { method: 'POST', auth: false, data: { phone } })
      .then(() => {
        this.setData({ counting: 60 });
        const timer = setInterval(() => {
          const n = this.data.counting - 1;
          this.setData({ counting: n });
          if (n <= 0) clearInterval(timer);
        }, 1000);
        toast('验证码已发送');
      })
      .catch(err => toast(err.message));
  },

  submit() {
    const { id, dates, dateIdx, slots, slotIdx, name, phone, code, submitting } = this.data;
    if (submitting) return;
    if (slotIdx < 0) return toast('先选个时间');
    if (!name.trim()) return toast('姓名必填');
    if (!/^1\d{10}$/.test(phone)) return toast('手机号格式不对');
    if (!/^\d{6}$/.test(code)) return toast('验证码是 6 位数字');

    const slot = slots[slotIdx];
    this.setData({ submitting: true });
    request('/public/bookings', {
      method: 'POST',
      auth: false,
      data: {
        coach_id: Number(id),
        date: dates[dateIdx].date,
        start: slot.start, end: slot.end,
        name: name.trim(), phone, sms_code: code,
        idem_key: uuid(),
      },
    })
      .then(() => {
        this.setData({ submitting: false, done: true });
        this.askSubscribe();
      })
      .catch(err => {
        this.setData({ submitting: false });
        toast(err.message);
      });
  },

  // 下单时一次性申请订阅消息
  askSubscribe() {
    const ids = Object.values(TPL_IDS).filter(Boolean);
    if (!ids.length || !wx.requestSubscribeMessage) return;
    wx.requestSubscribeMessage({ tmplIds: ids });
  },
});
