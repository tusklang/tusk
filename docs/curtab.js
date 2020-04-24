$(document).ready(() => {
  var path = window.location.pathname;

  //if it is / then it must be /index.php
  if (path == '/') path = '/index.php';

  $(`a.item[href=\'${path}\']`).addClass('active');
});
