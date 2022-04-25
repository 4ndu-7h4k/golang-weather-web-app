let suggestions;
let PORT = document.getElementById("portValue").value;
let DOMAIN = document.getElementById("domainValue").value;

async function getCityList() {
	cityList = await fetch(
		"http://"+DOMAIN+":"+PORT+"/getCity"
	).then((response) => {
		if (!response.ok) {
			console.log("Unalble to fecth time City");
			return;
		}
		return response.json();
	});
	const {
		Cities
	} = cityList;
	suggestions = Cities;
}
cityList = getCityList();

