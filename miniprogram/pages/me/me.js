const { request } = require('../../utils/request');
const { WEEKDAYS, fen2yuan, toast } = require('../../utils/util');

Page({
  data: {
    nickname: '', title: '', bio: '', home_court: '',
    priceYuan: '', video_url: '', pay_qr_url: '',
    published: false,
    days: [], // [{k:'0', name:'周日', text:'09:00-10:00,14:00-15:00'}]
    isAdmin: false,
  },

  onShow() {
    this.load();
    request('/admin/ping').then(() => this.setData({ isAdmin: true })).catch(() => {});
  },

  goAdmin() {
    wx.navigateTo({ url: '/pages/admin/admin' });
  },

  goClaim() {
    wx.navigateTo({ url: '/pages/claim/claim' });
  },

  load() {
    return request('/coach/me')
      .then(c => {
        let weekly = {};
        try { weekly = JSON.parse(c.schedule || '{}').weekly || {}; } catch (e) {}
        const days = [];
        for (let k = 0; k < 7; k++) {
          const slots = weekly[String(k)] || [];
          days.push({ k: String(k), name: WEEKDAYS[k], text: slots.map(s => `${s[0]}-${s[1]}`).join(',') });
        }
        this.setData({
          nickname: c.nickname || '', title: c.title || '', bio: c.bio || '',
          home_court: c.home_court || '', priceYuan: fen2yuan(c.price_per_hour),
          video_url: c.video_url || '', pay_qr_url: c.pay_qr_url || '',
          published: !!c.published, days,
        });
      })
      .catch(err => toast(err.message));
  },

  onInput(e) {
    this.setData({ [e.currentTarget.dataset.k]: e.detail.value });
  },

  onDayInput(e) {
    const k = e.currentTarget.dataset.k;
    const days = this.data.days.map(d => d.k === k ? { ...d, text: e.detail.value } : d);
    this.setData({ days });
  },

  // "09:00-10:00,14:00-15:00" → [["09:00","10:00"],...]
  parseDay(text) {
    const out = [];
    for (const part of String(text).split(/[,，、\s]+/)) {
      const t = part.trim();
      if (!t) continue;
      const m = t.match(/^(\d{1,2}:\d{2})-(\d{1,2}:\d{2})$/);
      if (!m || m[1] >= m[2]) throw new Error(`时段格式不对：${t}`);
      out.push([m[1].padStart(5, '0'), m[2].padStart(5, '0')]);
    }
    return out;
  },

  collectProfile() {
    const weekly = {};
    for (const d of this.data.days) weekly[d.k] = this.parseDay(d.text);
    const price = Math.round(parseFloat(this.data.priceYuan || '0') * 100);
    const { nickname, title, bio, home_court, video_url, pay_qr_url } = this.data;
    return {
      nickname, title, bio, home_court, video_url, pay_qr_url,
      price_per_hour: price,
      schedule: JSON.stringify({ weekly }),
    };
  },

  onSave() {
    let data;
    try {
      data = this.collectProfile();
    } catch (err) {
      return toast(err.message);
    }
    request('/coach/profile', { method: 'PUT', data })
      .then(() => toast('已保存'))
      .catch(err => toast(err.message));
  },

  onPublish() {
    const published = !this.data.published;
    wx.showModal({
      title: published ? '发布主页？' : '下架主页？',
      content: published ? '发布后学员可搜索到你' : '下架后学员看不到你的主页',
      success: res => {
        if (!res.confirm) return;
        // 先保存，再发布
        this.onSaveSilent(() => {
          request('/coach/publish', { method: 'POST', data: { published } })
            .then(() => {
              this.setData({ published });
              toast(published ? '已发布' : '已下架');
            })
            .catch(err => toast(err.message));
        });
      },
    });
  },

  onSaveSilent(cb) {
    let data;
    try {
      data = this.collectProfile();
    } catch (err) {
      return toast(err.message);
    }
    request('/coach/profile', { method: 'PUT', data }).then(cb).catch(err => toast(err.message));
  },
});
