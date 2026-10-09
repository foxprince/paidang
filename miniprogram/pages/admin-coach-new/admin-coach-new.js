const { request } = require('../../utils/request');
const { toast } = require('../../utils/util');

Page({
  data: { name: '', phone: '', title: '', priceYuan: '', bio: '', done: null },

  onInput(e) {
    this.setData({ [e.currentTarget.dataset.k]: e.detail.value });
  },

  onSubmit() {
    const { name, phone, title, priceYuan, bio } = this.data;
    if (!name.trim()) return toast('姓名必填');
    request('/admin/coaches', {
      method: 'POST',
      data: { name: name.trim(), phone, title, bio, price_yuan: parseFloat(priceYuan || '0') },
    })
      .then(coach => this.setData({ done: coach }))
      .catch(err => toast(err.message));
  },

  onCopy() {
    wx.setClipboardData({ data: this.data.done.claim_code || '' });
  },

  onAgain() {
    this.setData({ name: '', phone: '', title: '', priceYuan: '', bio: '', done: null });
  },
});
