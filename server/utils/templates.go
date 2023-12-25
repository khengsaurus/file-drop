package utils

const ImagePageHtml = `
<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="UTF-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<title>{{.Title}}</title>
	</head>
	<body style="margin: 0px">
		<img
			src="{{.Src}}"
			alt="Image"
			style="max-width: 100vw"
		/>
	</body>
</html>
`
