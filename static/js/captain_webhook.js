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
			return $scope.transformers.activeTransformer.Name;
		}
		return 'Configuration';
	}
}]);

captainWebhookApp.controller('captainWebhookTransformerEdit', ['$routeParams', '$scope', 'Transformers', function($routeParams, $scope, Transformers){
	$scope.transformer;
	$scope.messages;	
	$scope.transformerId = $routeParams.id;
	$scope.activeMessage;
	$scope.activeRelId;
	$scope.sortableRelIdList;

	Transformers.getOne($scope.transformerId, function(transformer){
		$scope.transformer = transformer;
		Transformers.activeTransformer = transformer;
	});

	Transformers.getMessages($scope.transformerId, function(messages){
		$scope.messages = messages;
		$scope.activeMessage = messages[0];
	});

	$scope.saveTransformer = function(){
		Transformers.saveTransformer($scope.transformer);
	}

	$scope.setActiveMessageObject = function(relId){
		$scope.activeRelId = relId;
		if(angular.isUndefined($scope.transformer.ObjectTransformation[relId])){
			$scope.transformer.ObjectTransformation[relId] = {
				Template : '',
				Rel_id : relId
			};
		}
	}

	$scope.addProperty = function(relId, key){
		$scope.transformer.ObjectTransformation[relId].Template = $scope.transformer.ObjectTransformation[relId].Template + " {{." + key + "}}"
	}

	$scope.setActiveMessage = function(message){		
		$scope.activeMessage = message;
	}
}]);