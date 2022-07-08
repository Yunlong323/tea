name: string
kind: string
#Classification:{
	[...Entity]
}

#Entity:{
	type: string
	task:{
		args:[
			...{name:string,type:string}
		],
		cmd: [...string]
	},
	#Service,
	#Nodes
}

#Service:{
		entry:[...string],
		api:[...{
			name:string,
			url:string,
			args:[...{
				name:string,
				type:string
			}],
			return:[...{name:string,type:string}]
		}]
}
#Nodes:{
	[...{
		name:string,
		socket:string
	}]
}






