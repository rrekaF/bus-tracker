        // Initialize and add the map
        let map;
        let chosenRoute;
        const form = document.getElementById("form");
        console.log('dupa')
        form.addEventListener("submit", submitForm);
        function submitForm(){
           chosenRoute = document.getElementById("busNumber").value 
           console.log(chosenRoute);
           getData(chosenRoute);
           event.preventDefault();

        }

        function getData(route) {
            let myBuses = [];
            let data = [];
            axios.get("https://ckan2.multimediagdansk.pl/gpsPositions?v=2").then(function (response) {
                for (const key in response['data']['vehicles']) {
                    let item = response['data']['vehicles'][key];
                    if (item['routeShortName'] == route) {
                        myBuses.push(item);
                    }
                }
            }).catch(function (error) {
                console.log("Error: " + error);
            }).finally(function () {
                // console.log(myBuses);
                for (bus in myBuses) {
                    console.log(bus)
                    data.push({
                        position: {
                            lat: myBuses[bus]['lat'],
                            lng: myBuses[bus]['lon']
                        },
                        tripId: myBuses[bus]['tripId'],
                        angle: myBuses[bus]['direction'],
                        headsign: myBuses[bus]['headsign']
                    });
                }
                console.log(data)
                initMap(data);
            })

        }

        async function initMap(data) {
            if(data.length == 0){
                map = document.getElementById("noBuses");
                map.innerHTML = "<h1>No buses " + chosenRoute + " at this time.</h1>"
                return;
            } else {
                map = document.getElementById("noBuses");
                map.innerHTML = ""
            }
            // Request needed libraries.
            const { Map } = await google.maps.importLibrary("maps");
            const { AdvancedMarkerElement } = await google.maps.importLibrary("marker");

            let mapCenter = { lat: 0, lng: 0 };
            for (bus in data) {
                mapCenter['lat'] += data[bus]['position']['lat'];
                mapCenter['lng'] += data[bus]['position']['lng'];
            }
            mapCenter['lat'] /= data.length;
            mapCenter['lng'] /= data.length;

            map = new Map(document.getElementById("map"), {
                zoom: 14,
                center: mapCenter,
                mapId: "demo",
            });

            let markers = [];
            let arrows = [];
            const lineSymbol = {
                path: google.maps.SymbolPath.FORWARD_CLOSED_ARROW,
            };

            for (bus in data) {
                markers.push(new AdvancedMarkerElement({
                    map: map,
                    position: data[bus]['position'],
                    title: data[bus]['headsign']
                }))

                let destLat, destLng;
                destLat = 0.00002 * data[bus]['position']['lat'] * Math.cos(data[bus]['angle'] * (Math.PI / 180));
                destLng = 0.0002 * data[bus]['position']['lng'] * Math.sin(data[bus]['angle'] * (Math.PI / 180));
                destLat += data[bus]['position']['lat'];
                destLng += data[bus]['position']['lng'];

                arrows.push(new google.maps.Polyline({
                    path: [
                        data[bus]['position'],
                        { lat: destLat, lng: destLng },
                    ],
                    icons: [
                        {
                            icon: lineSymbol,
                            offset: "100%",
                        },
                    ],
                    map: map,
                }))
            }
        }
