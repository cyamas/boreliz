function addEdgesToForm(numEdges) {

	console.log("NUM EDGES: " + numEdges)
	let edgeDiv = document.getElementById("edges-container")
	edgeDiv.innerHTML = ""

	let breakElem = document.createElement("br")
	for (let i = 0; i < numEdges; i++) {
		let label = document.createElement("div")
		label.innerHTML = `Edge ${i + 1}`
		edgeDiv.appendChild(label)

		let angleLabel = document.createElement("label")
		angleLabel.innerHTML = "Angle:"
		edgeDiv.appendChild(angleLabel)

		let angle = document.createElement("input")
		angle.id = `edge-${i}-angle`
		angle.name = `edge-${i}-angle`
		angle.type = "number"
		angle.min = "-179"
		angle.max = "180"
		angle.step = "1"
		angle.placeholder = "Angle"
		edgeDiv.appendChild(angle)
		edgeDiv.appendChild(breakElem.cloneNode(true))

		let widthLabel = document.createElement("label")
		widthLabel.innerHTML = "Width:"
		edgeDiv.appendChild(widthLabel)

		let width = document.createElement("input")
		width.id = `edge-${i}-width`
		width.name = `edge-${i}-width`
		width.type = "number"
		width.min = "1"
		width.max = "8"
		width.step = "1"
		width.placeholder = "Width"
		edgeDiv.appendChild(width)
		edgeDiv.appendChild(breakElem.cloneNode(true))

		let depthLabel = document.createElement("label")
		depthLabel.innerHTML = "Depth:"
		edgeDiv.appendChild(depthLabel)

		let depth = document.createElement("input")
		depth.id = `edge-${i}-depth`
		depth.name = `edge-${i}-depth`
		depth.type = "number"
		depth.min = "0.125"
		depth.max = "3.000"
		depth.step = "0.125"
		depth.placeholder = "Depth"
		edgeDiv.appendChild(depth)
	}
}
function zoom(dir) {
	let widthInput = document.getElementById("image-width-input")
	let heightInput = document.getElementById("image-height-input")
	let image = document.getElementById("hold-image")
	let currWidth = image.clientWidth;
	let currHeight = image.clientHeight;

	let diff = 2
	if (dir < 0) {
		diff = -2;
	}
	let newWidth = currWidth + diff;
	let newHeight = currHeight + diff;
	if (newWidth == 0 || newHeight == 0) {
		return
	}
	image.style.width = newWidth + "px";
	image.style.height = newHeight + "px";
	widthInput.value = newWidth;
	heightInput.value = newHeight;
}
function rotateHoldImage(angle) {
	let slider = document.getElementById("image-angle-output");
	let image = document.getElementById("hold-image");
	image.style.transform = `rotate(${angle}deg)`;
	slider.innerHTML = angle;
}

function addDragAndDrop() {
	console.log("addDragAndDrop() called");
	let dropArea = document.getElementById('hold-image-container');
	let image = document.createElement('img')
	image.id = "hold-image"
	let targetBolt = document.getElementsByClassName('img-bolt')[4]
	let fileInput = document.getElementById('form-image')

	dropArea.addEventListener('click', () => fileInput.click());
	console.log("click added to dropArea");

	// Handle drag and drop
	['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
		dropArea.addEventListener(eventName, preventDefaults, false);
	});

	function preventDefaults(e) {
		e.preventDefault();
		e.stopPropagation();
	}

	['dragenter', 'dragover'].forEach(eventName => {
		dropArea.addEventListener(eventName, () => dropArea.classList.add('drag-over'), false);
	});

	['dragleave', 'drop'].forEach(eventName => {
		dropArea.addEventListener(eventName, () => dropArea.classList.remove('drag-over'), false);
	});

	dropArea.addEventListener('drop', handleDrop, false);

	function handleDrop(e) {
		console.log("file dropped")
		let dt = e.dataTransfer;
		let files = dt.files;
		if (files.length > 1) {
			showError("You can only add one image!");
			return;
		}

		let reader = new FileReader();
		reader.onload = function(e) {
			const img = new Image();
			img.onload = function() {
				if (img.width !== img.height) {
					showError("Error: Image must have square dimensions")
					return
				}
				image.src = e.target.result;
				image.style.width = '24px'
				image.style.height = '24px'
				targetBolt.innerHTML = '';
				targetBolt.appendChild(image);
			};
			img.src = e.target.result;
		}

		reader.readAsDataURL(files[0])
		fileInput.files = files;
	}
}

function showError(message) {
	alert(message);
}
