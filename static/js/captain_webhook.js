var captainWebhookApp = angular.module('captainWebhookApp', ['ui.bootstrap', 'ngRoute'])
	.config(['$routeProvider', '$locationProvider', 
		function($routeProvider, $locationProvider){
			$routeProvider
				.when('/transformer/:id', {
		          templateUrl: '/html/transformer_edit.html',
		          controller: 'captainWebhookTransformerEdit',
		        });
		    $locationProvider.html5Mode(true);
		}]);

captainWebhookApp.controller('captainWebhook', ['$scope', 'Transformers', '$route', '$routeParams', '$location', function($scope, Transformers, $route, $routeParams, $location){
	$scope.transformers = Transformers;
	this.$route = $route;
    this.$location = $location;
    this.$routeParams = $routeParams;
	Transformers.get();

	$scope.activeTransformerName = function(){
		if(angular.isDefined($scope.transformers.activeTransformer)){
			return $scope.transformers.activeTransformer.name;
		}
		return 'Configuration';
	}
}]);

captainWebhookApp.controller('captainWebhookTransformerEdit', ['$routeParams', '$scope', 'Transformers', function($routeParams, $scope, Transformers){
	$scope.transformer;
	$scope.messages;	
	$scope.transformerId = $routeParams.id;
	$scope.activeMessageIndex;

	Transformers.getOne($scope.transformerId, function(transformer){
		$scope.transformer = transformer;
		Transformers.activeTransformer = transformer;
	});

	Transformers.getMessages($scope.transformerId, function(messages){
		$scope.messages = messages;
		$scope.activeMessageIndex = 0;
	});

	$scope.showMessage = function(messageIndex){
		$scope.activeMessageIndex = messageIndex;
	}
}]);