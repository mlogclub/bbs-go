import HttpClient from './HttpClient'

class CategoryApi {
  get(id) {
    return HttpClient.get("/api/admin/category/" + id)
  }

  list(params) {
    return HttpClient.post('/api/admin/category/list', params)
  }

  create(form) {
    return HttpClient.post('/api/admin/category/create', form)
  }

  update(form) {
    return HttpClient.post("/api/admin/category/update", form)
  }

  options() {
    return HttpClient.get("/api/admin/category/options")
  }
}

export default new CategoryApi()
