// 上线前改成正式域名，且必须 https
module.exports = {
  API_BASE: 'http://localhost:8080/api/v1',
  // 订阅消息模板 ID（与服务端 .env 里的三个保持一致，申请后填写）
  TPL_IDS: {
    new_booking: '',
    result: '',
    done: '',
  },
};
