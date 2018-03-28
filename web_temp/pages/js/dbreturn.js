function sendArchive(arg) {
	document.getElementById("info").write(arg);
}

function updateAddImage(arg) {
	document.getElementById('image').value = arg;
}


window.addEventListener('DOMContentLoaded', function () {
	document.getElementById('updcanvas').addEventListener('click',function(){
		updateAddImage(document.getElementById('capture').value);
	}
}