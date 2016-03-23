
//	Actions

//	Stores

class CentralConfigAPIUtils {

	constructor(){

	}

	//	Gets all configuration items from the server
	getAllConfigItems(){

		//  Format the url
        let url = `https://query.yahooapis.com/v1/public/yql?q=select * from weather.forecast where woeid in (SELECT woeid FROM geo.placefinder WHERE text="${latitude},${longitude}" and gflags="R")&format=json`;

        $.ajax( url )
        .done(function(data) {
            //  Call the action to receive the data:
            //	WeatherActions.recieveWeatherData(weatherdata);
        }.bind(this))
        .fail(function() {
            //  Something bad happened
            console.log("There was a problem getting config items");
        });
	}
}