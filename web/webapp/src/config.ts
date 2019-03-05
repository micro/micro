export default {
  url: {
    basicUrl: process.env.NODE_ENV === "development" ? "/v1/admin" : "/v1/admin"
  }
};
