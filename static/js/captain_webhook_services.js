captainWebhookApp.service('Transformers', ["$http", "$rootScope", function($http, $rootScope){
	var Transformers = {};
	Transformers.transformers;
	Transformers.activeTransformer;
	Transformers.get = function(){
		if(typeof Transformers.transformers == 'undefined'){
			$http.get('/api/transformers').success(function(response){
				Transformers.transformers = response;
			});
		}
	}

	Transformers.getOne = function(id, callback){
		$http.get('/api/transformers/' + id).success(function(response){
			callback(response);
		});
	}

	Transformers.getMessages = function(id, callback){
		$http.get('/api/transformers/' + id + '/messages').success(function(response){
			callback(response);
		});
	}

	Transformers.saveTransformer = function(transformer){
		$http.put('/api/transformers/' + transformer.Id, transformer).success(function(response){
			
		})
	}

	return Transformers;
}]);