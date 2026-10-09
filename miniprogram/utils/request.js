const { API_BASE } = require('./config');

// 统一请求：自动带 token，code != 0 时抛错
function request(path, { method = 'GET', data = {}, auth = true } = {}) {
  return new Promise((resolve, reject) => {
    const header = { 'content-type': 'application/json' };
    if (auth) {
      const token = wx.getStorageSync('token');
      if (token) header['Authorization'] = 'Bearer ' + token;
    }
    wx.request({
      url: API_BASE + path,
      method,
      data,
      header,
      success(res) {
        const body = res.data || {};
        if (body.code === 0) {
          resolve(body.data);
        } else if (body.code === 40101) {
          wx.removeStorageSync('token');
          reject(new Error('登录已过期'));
        } else {
          reject(new Error(body.msg || '请求失败'));
        }
      },
      fail(err) {
        reject(err);
      },
    });
  });
}

module.exports = { request };
