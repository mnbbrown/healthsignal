import config from "./config";

const getUrl = path => {
  if (path[0] !== "/") path = `/${path}`;
  return config.api + path;
};

const api = {
  getEndpoints: () => {
    return fetch(getUrl("endpoints")).then(response => response.json());
  },
  getEndpointData: endpoint => {
    return fetch(getUrl(`/endpoints/${endpoint}/data`)).then(response =>
      response.json()
    );
  }
};

export default api;
