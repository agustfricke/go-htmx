<script src="public/htmx.min.js"></script>;
document.addEventListener("htmx:afterSwap", function (event) {
  if (event.target.id === "task-list") {
    document.getElementById("name").value = "";
  }
});
