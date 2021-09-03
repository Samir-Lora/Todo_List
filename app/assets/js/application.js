require("jquery");
require("bootstrap/dist/js/bootstrap.bundle.js");
require("@fortawesome/fontawesome-free/js/all.js");

$(() => {
  $("body").on("click", '#idb', function () {
    var r = '<a href="/task/delete/' + $(this).attr('data-id') + '" data-method="DELETE" ><button type="button" class="btn btn-primary">Delete</button></a>';
    document.getElementById('deleteb').innerHTML = r;
  });
  $("body").on("click", '#modaluser', function () {
    var r = '<a href="/user/delete/' + $(this).attr('data-id') + '" data-method="DELETE" ><button type="button" class="btn btn-primary">Delete</button></a>';
    document.getElementById('deletebtn').innerHTML = r;
  });
  $("body").on("click", '#modalactive', function () {
    var r = '<a href="user/updateactive/' + $(this).attr('data-id') + '" data-method="PUT" ><button type="button" class="btn btn-primary">Deactivate</button></a>';
    document.getElementById('inactivebtn').innerHTML = r;
  });
});
