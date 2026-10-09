const { request } = require('../../utils/request');
const { toast } = require('../../utils/util');

Page({
  data: {
    id: null,
    isNew: false,
    student: null,
    notes: [],
    // 新建表单
    name: '', phone: '', wechat: '', level: '',
    // 备注表单
    keyPoints: '', nextPlan: '',
  },

  onLoad(options) {
    const id = options.id;
    if (id === 'new') {
      this.setData({ isNew: true });
      wx.setNavigationBarTitle({ title: '新增学员' });
    } else {
      this.setData({ id });
      this.load();
    }
  },

  load() {
    return request(`/students/${this.data.id}`)
      .then(({ student, notes }) => this.setData({
        student,
        notes: (notes || []).map(n => ({ ...n, dateLabel: (n.created_at || '').slice(0, 10) })),
      }))
      .catch(err => toast(err.message));
  },

  onInput(e) {
    this.setData({ [e.currentTarget.dataset.k]: e.detail.value });
  },

  onCreate() {
    const { name, phone, wechat, level } = this.data;
    if (!name.trim()) return toast('姓名必填');
    request('/students', { method: 'POST', data: { name, phone, wechat, level } })
      .then(() => wx.navigateBack())
      .catch(err => toast(err.message));
  },

  onAddNote() {
    const { id, keyPoints, nextPlan } = this.data;
    if (!keyPoints.trim()) return toast('本节要点不能为空');
    request(`/students/${id}/notes`, { method: 'POST', data: { key_points: keyPoints, next_plan: nextPlan } })
      .then(() => this.setData({ keyPoints: '', nextPlan: '' }))
      .then(() => this.load())
      .catch(err => toast(err.message));
  },
});
