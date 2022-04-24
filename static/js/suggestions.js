let suggestions;
async function getCityList() {
	cityList = await fetch(
		"http://localhost:8080/getCity"
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

