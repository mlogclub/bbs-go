import HttpClient from './HttpClient'

class TagApi {
  get(id) {
    return HttpClient.get("/api/admin/tag/" + id)
  }

  list() {
    return HttpClient.post('/api/admin/tag/list')
  }

  create(form) {
    return HttpClient.post('/api/admin/tag/create', form)
  }

  update(form) {
    return HttpClient.post("/api/admin/tag/update", form)
  }

}

export default new TagApi()
