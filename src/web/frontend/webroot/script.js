var app = angular.module('myApp', []);
app.controller('customersCtrl', function($scope, $http) {
  //get data from API
  $http.get("/api/list")
  .then(function (response) {
    $scope.names = response.data;

    angular.forEach($scope.names, function (item) {
      if (item.Enabled){
        item.selectedOption = {id: '1', name: 'true'}
      } else{
        item.selectedOption = {id: '2', name: 'false'}
      }
    })
  });

  //onclick handler
  $scope.remove = function(array, index){
    array.splice(index, 1);
  }

  $scope.update = function(array, index, id){
    var name = $('#' + id + '_name').val()
    var enabled = $('#' + id + '_enabled').val()

    // i can just pass the index and the array, so i dont need to do a for loop ^^
    array[index].Name = name;
    array[index].Enabled = enabled;

    console.log(typeof(array[index]))

    //Call the services
    $http.post('/api/update', JSON.stringify(array[index])).then(function (response) {
      if (response.data)
        $scope.msg = "Post Data Submitted Successfully!";
      }, function (response) {
        $scope.msg = "Service not Exists";
        $scope.statusval = response.status;
        $scope.statustext = response.statusText;
        $scope.headers = response.headers();
    });
  }

  $scope.data = {
    availableOptions: [
      {id: '1', name: 'true'},
      {id: '2', name: 'false'}
    ]
  };
});
