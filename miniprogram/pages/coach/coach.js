const { request } = require('../../utils/request');
const { fen2yuan, toast } = require('../../utils/util');

Page({
  data: { id: null, coach: null, avgRating: '0.0', reviews: [], loading: true },

  onLoad(options) {
    this.setData({ id: options.id });
    this.load();
  },

  load() {
    return request(`/public/coach/${this.data.id}`, { auth: false })
      .then(({ coach, avg_rating, reviews }) => this.setData({
        coach: { ...coach, priceYuan: fen2yuan(coach.price_per_hour) },
        avgRating: (avg_rating || 0).toFixed(1),
        reviews: (reviews || []).map(r => ({ ...r, stars: '★'.repeat(r.rating) + '☆'.repeat(5 - r.rating) })),
        loading: false,
      }))
      .catch(err => {
        this.setData({ loading: false });
        toast(err.message);
      });
  },

  goBook() {
    wx.navigateTo({ url: `/pages/book/book?id=${this.data.id}` });
  },

  onShareAppMessage() {
    const { coach } = this.data;
    return { title: `${coach.nickname} · 网球陪练`, path: `/pages/coach/coach?id=${this.data.id}` };
  },
});
