interface Options {
  host: string
}

export class Api {
  options: Options

  constructor(options: Options) {
    this.options = options
  }

  async storeDefaultConfig(category: string, name: string, data: object) {
    let url = `/api/v1/default-configs/${category}/${name}`

    return this.request('POST', url, data)
  }

  async getDefaultConfigByName(category: string, name: string) {
    let url = `/api/v1/default-configs/${category}/${name}`

    return this.request('GET', url)
  }

  async getDefaultConfigs(category: string) {
    return this.request('GET', '/api/v1/default-configs/statics')
  }

  async storeClient(id: string, name: string, secret: string) {
    return this.request('POST', '/api/v1/clients', {
      id,
      name,
      secret,
    })
  }

  async deleteClient(clientId: string) {
    return this.request('DELETE', `/api/v1/clients/${clientId}`)
  }

  async clientList(page: number, limit: number = 100) {
    return this.request('GET', `/api/v1/clients?page=${page}&limit=${limit}`)
  }

  async storeClientConfig(clientId: string, category: string, name: string, parent: string, data: object) {
    return this.request('POST', `/api/v1/clients/${clientId}/${category}/${name}`, {
      parent: parent,
      data: data,
    })
  }

  async deleteClientConfig(clientId: string, name: string) {
    return this.request('DELETE', `/api/v1/clients/${clientId}/static/${name}`)
  }

  async listClientConfigStatic(clientId: string) {
    return this.request('GET', `/api/v1/clients/${clientId}/static`)
  }

  async getClientConfigStaticById(clientId: string, id: string) {
    return this.request('GET', `/api/v1/clients/${clientId}/static/${id}`)
  }

  async deleteStaticElement(clientId: string, staticId: string, staticKey: string) {
    return this.request('DELETE', `/api/v1/clients/${clientId}/static/${staticId}/${staticKey}`)
  }

  async deleteDefaultStaticElement(staticId: string, staticKey: string) {
    return this.request('DELETE', `/api/v1/default-configs/statics/${staticId}/${staticKey}`)
  }

  private async request(method: string, url: string, data?: any) {
    url = this.options.host.replace(/\/$/, '') + ('/' + url).replace('//', '/')

    let options: RequestInit = {
      method: method,
    }

    if (data) {
      options.body = JSON.stringify(data)
    }

    return await fetch(url, options)
  }
}
