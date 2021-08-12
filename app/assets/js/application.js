require("jquery");
require("bootstrap/dist/js/bootstrap.bundle.js");
require("@fortawesome/fontawesome-free/js/all.js");

$(() => {
  $('body').on('click', '#idn1', function () {
    location.href = "/task/new";
  })
  $("body").on("click", '#idb', function () {
    var r = '<a href="/task/delete/' + $(this).attr('class') + '" data-method="DELETE" ><button type="button" class="btn btn-primary">Delete</button></a>';
    document.getElementById('deleteb').innerHTML = r;
  });
  function currentDate() {
    var week = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
    var month = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]
    var d = new Date();
    var dayWeek = week[d.getDay()];
    var dayMonth = d.getDate();
    var Month = month[d.getMonth() + 1];
    var year = d.getFullYear();

    document.getElementById("date").innerHTML = `${dayWeek} ${dayMonth} , ${Month} ${year}`;
  }
  currentDate();
  var URLsearch = window.location.search;
  if (URLsearch ==  "?complete=true") {
    $('.addbtn').addClass('d');
    a = document.getElementById('com').style.fontWeight = "bold";
  } else if (URLsearch ==  "?complete=false"){
    a = document.getElementById('textincompleted').style.fontWeight = "bold";

  }
});
