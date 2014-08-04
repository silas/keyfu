delete this.console;

this.link = function(url, queryUrl) {
  if (query) {
    if (!queryUrl) {
      throw new Error('no query url found');
    }
    this.location = queryUrl.replace(/%s/g, encodeURIComponent(query));
  } else if (url) {
    this.location = url;
  } else {
    throw new Error('no url found');
  }
};
