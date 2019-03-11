# Image Processing API documentation:

Path: /process
Method: HTTP Post
Description: Execute one or more image transformation operations, in-order, on the specified image.

Request Body Parameters:
url: The HTTP URL of the image to transform
operations: a list of one or more image transformation operations, and each operation is specified as a string which matches the supported vocabulary which is described below.

Request body example:
{
"url":"https://i.imgur.com/AOXD2P4.jpg", 
"operations":[“flipVertical”, “rotateRight”, “resize,88”]
}

Operation vocabulary and usage:

flipVertical: This flips the image vertically.

flipHorizontal: This flips the image horizontally.

rotateRight: This rotates the image by 90 degrees clockwise.

rotateLeft: This rotates the image by 90 degrees counter-clockwise.

rotate: This is the same as rotateRight.
rotate,n: This rotates the image by n degrees where n is an integer.

grayscale: This converts the image into grayscale.

resize: This defaults to resizing to a 100 pixel wide image while preserving the aspect ratio.
resize,n: This resizes the image to an n-pixel wide image while preserving the aspect ratio.
resize,w,h: This resizes the image to a w-by-h pixel image, which could result in stretching the image in one or both directions.

thumbnail: This defaults to resizing to a 100 pixel wide image while preserving the aspect ratio.
