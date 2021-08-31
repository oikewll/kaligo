import 'dart:convert';
import 'dart:mirrors';
import 'package:http/http.dart' as http;

final future = http.get("https://flutterdevs.com");
future.then((response) {
  if (response.statusCode == 200) {
    print("Response received.");
  }
});

