export default [
  {
    name: "APP_LOGIN_SUCCESS",
    callback: function(e) {
      this.$router.push({
        path: "dashboard"
      });
    }
  },
  {
    name: "APP_PAGE_LOADED",
    callback: function(e) {}
  },
  {
    name: "APP_AUTH_FAILED",
    callback: function(e) {
      this.$router.push("/login");
      this.$message.error("Token has expired");
    }
  },
  {
    name: "APP_BAD_REQUEST",
    // @error api response data
    callback: function(msg) {
      this.$message.error(msg);
    }
  },
  {
    name: "APP_ACCESS_DENIED",
    // @error api response data
    callback: function(msg) {
      this.$message.error(msg);
      this.$router.push("/forbidden");
    }
  },
  {
    name: "APP_RESOURCE_DELETED",
    // @error api response data
    callback: function(msg) {
      this.$message.success(msg);
    }
  },
  {
    name: "APP_RESOURCE_UPDATED",
    // @error api response data
    callback: function(msg) {
      this.$message.success(msg);
    }
  }
];
