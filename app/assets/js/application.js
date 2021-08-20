require("jquery");
require("bootstrap/dist/js/bootstrap.bundle.js");
require("@fortawesome/fontawesome-free/js/all.js");

$(() => {
  $("body").on("click", '#idb', function () {
    var r = '<a href="/task/delete/' + $(this).attr('data-id') + '" data-method="DELETE" ><button type="button" class="btn btn-primary">Delete</button></a>';
    document.getElementById('deleteb').innerHTML = r;
  });
});
