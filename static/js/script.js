const searchWrapper = document.querySelector(".search");
const inputBox = document.getElementById("sbox");
const sugBox = document.querySelector(".autocom-box");
const icon = document.querySelector(".search-icon");
let weather = {
	fetchWeather: function (city) {
		fetch(
				"http://localhost:8080/weather/" +
				city
			)
			.then((response) => {
				if (!response.ok) {
					alert("No weather found.");
					throw new Error("No weather found.");
				}
				return response.json();
			})
			.then((data) => this.displayWeather(data));
	},
	displayWeather: async function (data) {
		const {	name	} = data;
		const {
			country
		} = data.sys;
		const {
			icon,
			description
		} = data.weather[0];
		const {
			temp,
			humidity
		} = data.main;
		const {
			timezone
		} = data;
		const {
			speed
		} = data.wind;
		locTime = await fetch(
			"http://localhost:8080/time/" +
			timezone
		).then((response) => {
			if (!response.ok) {
				console.log("Unalble to fecth time api");
				return;

			}
			return response.json();
		});
		const {
			time
		} = locTime;
		document.querySelector(".city").innerText = name + "," + country;
		document.querySelector(".time").innerText = time;
		document.querySelector(".icon").src =
			"https://openweathermap.org/img/wn/" + icon + ".png";
		document.querySelector(".description").innerText = description;
		document.querySelector(".temp").innerText = Math.round(temp) + "°";
		document.querySelector(".humidity").innerText =
			"Humidity: " + humidity + "%";
		document.querySelector(".wind").innerText =
			"Wind speed: " + speed + " km/h";
		document.querySelector(".weather").classList.remove("loading");

	},
	search: function () {
		searchWrapper.classList.remove("active"); //hide autocomplete box
		console.log("Seraching Weather in city " + document.getElementById("sbox").value);
		this.fetchWeather(document.querySelector(".search-bar").value);
	},
};

document.querySelector(".search button").addEventListener("click", function () {
	weather.search();
});

document
	.querySelector(".search-bar")
	.addEventListener("keyup", function (event) {
		if (event.key == "Enter") {
			weather.search();
		}
	});
inputBox.onkeyup = (e) => {
	let userInputData = e.target.value;
	let emptyArray = [];
	if (userInputData) {
		emptyArray = suggestions.filter((data) => {
			return data.toLocaleLowerCase().startsWith(userInputData.toLocaleLowerCase());

		});

		emptyArray = emptyArray.map((data) => {
			return data = '<li>' + data + '</li>';
		});
		searchWrapper.classList.add("active"); //show autocomplete box
		showSuggestions(emptyArray);
		let allList = sugBox.querySelectorAll("li");
		for (let i = 0; i < allList.length; i++) {
			//adding onclick attribute in all li tag
			allList[i].setAttribute("onclick", "select(this)");
		}
	} else {
		searchWrapper.classList.remove("active"); //hide autocomplete box
	}

}

function select(element) {
	let selectData = element.textContent;
	inputBox.value = selectData;
	weather.fetchWeather(selectData);
	searchWrapper.classList.remove("active");
}

function showSuggestions(list) {
	let listData;
	if (!list.length) {
		userValue = inputBox.value;
		listData = `<li>${userValue}</li>`;
	} else {
		listData = list.join('');
	}
	sugBox.innerHTML = listData;
}