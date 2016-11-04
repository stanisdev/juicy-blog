document.addEventListener("DOMContentLoaded", function(event) { 
  var btn = document.getElementById('closeFlash');
  if (btn instanceof Object) {
    btn.onclick = function() {
      document.getElementById('flash').style.display = 'none';
    };
  }
});