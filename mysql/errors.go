package mysql

import (
	"fmt"
    "errors"
)

var (
	// ErrRecordNotFound record not found error
	ErrRecordNotFound = errors.New("record not found")
	// ErrInvalidTransaction invalid transaction when you are trying to `Commit` or `Rollback`
	ErrInvalidTransaction = errors.New("invalid transaction")
)

// ClientError is a type for mymysql client errors.
type ClientError string

func (e ClientError) Error() string {
	return string(e)
}

var (
    // ErrSeq packet sequence error
	ErrSeq            = ClientError("packet sequence error")
    // ErrPkt malformed packet
	ErrPkt            = ClientError("malformed packet")
    // ErrPktLong packet too long
	ErrPktLong        = ClientError("packet too long")
    // ErrUnexpNullLCS unexpected NULL LCS
	ErrUnexpNullLCS   = ClientError("unexpected NULL LCS")
    // ErrUnexpNullLCB unexpected NULL LCB
	ErrUnexpNullLCB   = ClientError("unexpected NULL LCB")
    // ErrUnexpNullDate unexpected NULL DATETIME
	ErrUnexpNullDate  = ClientError("unexpected NULL DATETIME")
    // ErrUnexpNullTime unexpected NULL TIME
	ErrUnexpNullTime  = ClientError("unexpected NULL TIME")
    // ErrUnkResultPkt unexpected or unknown result packet
	ErrUnkResultPkt   = ClientError("unexpected or unknown result packet")
    // ErrNotConn not connected
	ErrNotConn        = ClientError("not connected")
    // ErrAlredyConn already connected
	ErrAlredyConn     = ClientError("already connected")
    // ErrBadResult unexpected result
	ErrBadResult      = ClientError("unexpected result")
    // ErrUnreadedReply reply is not completely read
	ErrUnreadedReply  = ClientError("reply is not completely read")
    // ErrBindCount wrong number of values for bind
	ErrBindCount      = ClientError("wrong number of values for bind")
    // ErrBindUnkType unknown value type for bind
	ErrBindUnkType    = ClientError("unknown value type for bind")
    // ErrRowLength wrong length of row slice
	ErrRowLength      = ClientError("wrong length of row slice")
    // ErrBadCommand comand isn't text SQL nor *Stmt
	ErrBadCommand     = ClientError("comand isn't text SQL nor *Stmt")
    // ErrWrongDateLen wrong datetime/timestamp length
	ErrWrongDateLen   = ClientError("wrong datetime/timestamp length")
    // ErrWrongTimeLen wrong time length
	ErrWrongTimeLen   = ClientError("wrong time length")
    // ErrUnkMySQLType unknown MySQL type
	ErrUnkMySQLType   = ClientError("unknown MySQL type")
    // ErrWrongParamNum wrong parameter number
	ErrWrongParamNum  = ClientError("wrong parameter number")
    // ErrUnkDataType unknown data source type
	ErrUnkDataType    = ClientError("unknown data source type")
    // ErrSmallPktSize specified packet size is to small
	ErrSmallPktSize   = ClientError("specified packet size is to small")
    // ErrReadAfterEOR previous ScanRow call returned io.EOF
	ErrReadAfterEOR   = ClientError("previous ScanRow call returned io.EOF")
    // ErrOldProtocol server does not support 4.1 protocol
	ErrOldProtocol    = ClientError("server does not support 4.1 protocol")
    // ErrAuthentication authentication error
	ErrAuthentication = ClientError("authentication error")
)

// Error is a mymysql error.
// If a function/method returns error you may check the returned error type via
// a type assertion. If the type assertion succeeds you can retrieve the MySQL
// error code as below.
//
// Example:
//     if val, ok := err.(*mysql.Error); ok {
//         fmt.Println(val.Code)
//     }
type Error struct {
	Code uint16
	Msg  []byte
}

func (err Error) Error() string {
	return fmt.Sprintf("Received #%d error from MySQL server: \"%s\"",
		err.Code, err.Msg)
}

// MySQL error codes.
const (
	ErHashchk                                 = 1000
	ErNISAMCHK                                = 1001
	ErNO                                      = 1002
	ErYES                                     = 1003
	ErCantCreateFile                          = 1004
	ErCantCreateTable                         = 1005
	ErCantCreateDd                            = 1006
	//ErDB_CREATE_EXISTS                        = 1007
	//ErDB_DROP_EXISTS                          = 1008
	//ErDB_DROP_DELETE                          = 1009
	//ErDB_DROP_RMDIR                           = 1010
	//ErCANT_DELETE_FILE                        = 1011
	//ErCANT_FIND_SYSTEM_REC                    = 1012
	//ErCANT_GET_STAT                           = 1013
	//ErCANT_GET_WD                             = 1014
	//ErCANT_LOCK                               = 1015
	//ErCANT_OPEN_FILE                          = 1016
	//ErFILE_NOT_FOUND                          = 1017
	//ErCANT_READ_DIR                           = 1018
	//ErCANT_SET_WD                             = 1019
	//ErCHECKREAD                               = 1020
	//ErDISK_FULL                               = 1021
	//ErDUP_KEY                                 = 1022
	//ErERROR_ON_CLOSE                          = 1023
	//ErERROR_ON_READ                           = 1024
	//ErERROR_ON_RENAME                         = 1025
	//ErERROR_ON_WRITE                          = 1026
	//ErFILE_USED                               = 1027
	//ErFILSORT_ABORT                           = 1028
	//ErFORM_NOT_FOUND                          = 1029
	//ErGET_ERRNO                               = 1030
	//ErILLEGAL_HA                              = 1031
	//ErKEY_NOT_FOUND                           = 1032
	//ErNOT_FORM_FILE                           = 1033
	//ErNOT_KEYFILE                             = 1034
	//ErOLD_KEYFILE                             = 1035
	//ErOPEN_AS_READONLY                        = 1036
	//ErOUTOFMEMORY                             = 1037
	//ErOUT_OF_SORTMEMORY                       = 1038
	//ErUNEXPECTED_EOF                          = 1039
	//ErCON_COUNT_ERROR                         = 1040
	//ErOUT_OF_RESOURCES                        = 1041
	//ErBAD_HOST_ERROR                          = 1042
	//ErHANDSHAKE_ERROR                         = 1043
	//ErDBACCESS_DENIED_ERROR                   = 1044
	//ErACCESS_DENIED_ERROR                     = 1045
	//ErNO_DB_ERROR                             = 1046
	//ErUNKNOWN_COM_ERROR                       = 1047
	//ErBAD_NULL_ERROR                          = 1048
	//ErBAD_DB_ERROR                            = 1049
	//ErTABLE_EXISTS_ERROR                      = 1050
	//ErBAD_TABLE_ERROR                         = 1051
	//ErNON_UNIQ_ERROR                          = 1052
	//ErSERVER_SHUTDOWN                         = 1053
	//ErBAD_FIELD_ERROR                         = 1054
	//ErWRONG_FIELD_WITH_GROUP                  = 1055
	//ErWRONG_GROUP_FIELD                       = 1056
	//ErWRONG_SUM_SELECT                        = 1057
	//ErWRONG_VALUE_COUNT                       = 1058
	//ErTOO_LONG_IDENT                          = 1059
	//ErDUP_FIELDNAME                           = 1060
	//ErDUP_KEYNAME                             = 1061
	//ErDUP_ENTRY                               = 1062
	//ErWRONG_FIELD_SPEC                        = 1063
	//ErPARSE_ERROR                             = 1064
	//ErEMPTY_QUERY                             = 1065
	//ErNONUNIQ_TABLE                           = 1066
	//ErINVALID_DEFAULT                         = 1067
	//ErMULTIPLE_PRI_KEY                        = 1068
	//ErTOO_MANY_KEYS                           = 1069
	//ErTOO_MANY_KEY_PARTS                      = 1070
	//ErTOO_LONG_KEY                            = 1071
	//ErKEY_COLUMN_DOES_NOT_EXITS               = 1072
	//ErBLOB_USED_AS_KEY                        = 1073
	//ErTOO_BIG_FIELDLENGTH                     = 1074
	//ErWRONG_AUTO_KEY                          = 1075
	//ErREADY                                   = 1076
	//ErNORMAL_SHUTDOWN                         = 1077
	//ErGOT_SIGNAL                              = 1078
	//ErSHUTDOWN_COMPLETE                       = 1079
	//ErFORCING_CLOSE                           = 1080
	//ErIPSOCK_ERROR                            = 1081
	//ErNO_SUCH_INDEX                           = 1082
	//ErWRONG_FIELD_TERMINATORS                 = 1083
	//ErBLOBS_AND_NO_TERMINATED                 = 1084
	//ErTEXTFILE_NOT_READABLE                   = 1085
	//ErFILE_EXISTS_ERROR                       = 1086
	//ErLOAD_INFO                               = 1087
	//ErALTER_INFO                              = 1088
	//ErWRONG_SUB_KEY                           = 1089
	//ErCANT_REMOVE_ALL_FIELDS                  = 1090
	//ErCANT_DROP_FIELD_OR_KEY                  = 1091
	//ErINSERT_INFO                             = 1092
	//ErUPDATE_TABLE_USED                       = 1093
	//ErNO_SUCH_THREAD                          = 1094
	//ErKILL_DENIED_ERROR                       = 1095
	//ErNO_TABLES_USED                          = 1096
	//ErTOO_BIG_SET                             = 1097
	//ErNO_UNIQUE_LOGFILE                       = 1098
	//ErTABLE_NOT_LOCKED_FOR_WRITE              = 1099
	//ErTABLE_NOT_LOCKED                        = 1100
	//ErBLOB_CANT_HAVE_DEFAULT                  = 1101
	//ErWRONG_DB_NAME                           = 1102
	//ErWRONG_TABLE_NAME                        = 1103
	//ErTOO_BIG_SELECT                          = 1104
	//ErUNKNOWN_ERROR                           = 1105
	//ErUNKNOWN_PROCEDURE                       = 1106
	//ErWRONG_PARAMCOUNT_TO_PROCEDURE           = 1107
	//ErWRONG_PARAMETERS_TO_PROCEDURE           = 1108
	//ErUNKNOWN_TABLE                           = 1109
	//ErFIELD_SPECIFIED_TWICE                   = 1110
	//ErINVALID_GROUP_FUNC_USE                  = 1111
	//ErUNSUPPORTED_EXTENSION                   = 1112
	//ErTABLE_MUST_HAVE_COLUMNS                 = 1113
	//ErRECORD_FILE_FULL                        = 1114
	//ErUNKNOWN_CHARACTER_SET                   = 1115
	//ErTOO_MANY_TABLES                         = 1116
	//ErTOO_MANY_FIELDS                         = 1117
	//ErTOO_BIG_ROWSIZE                         = 1118
	//ErSTACK_OVERRUN                           = 1119
	//ErWRONG_OUTER_JOIN                        = 1120
	//ErNULL_COLUMN_IN_INDEX                    = 1121
	//ErCANT_FIND_UDF                           = 1122
	//ErCANT_INITIALIZE_UDF                     = 1123
	//ErUDF_NO_PATHS                            = 1124
	//ErUDF_EXISTS                              = 1125
	//ErCANT_OPEN_LIBRARY                       = 1126
	//ErCANT_FIND_DL_ENTRY                      = 1127
	//ErFUNCTION_NOT_DEFINED                    = 1128
	//ErHOST_IS_BLOCKED                         = 1129
	//ErHOST_NOT_PRIVILEGED                     = 1130
	//ErPASSWORD_ANONYMOUS_USER                 = 1131
	//ErPASSWORD_NOT_ALLOWED                    = 1132
	//ErPASSWORD_NO_MATCH                       = 1133
	//ErUPDATE_INFO                             = 1134
	//ErCANT_CREATE_THREAD                      = 1135
	//ErWRONG_VALUE_COUNT_ON_ROW                = 1136
	//ErCANT_REOPEN_TABLE                       = 1137
	//ErINVALID_USE_OF_NULL                     = 1138
	//ErREGEXP_ERROR                            = 1139
	//ErMIX_OF_GROUP_FUNC_AND_FIELDS            = 1140
	//ErNONEXISTING_GRANT                       = 1141
	//ErTABLEACCESS_DENIED_ERROR                = 1142
	//ErCOLUMNACCESS_DENIED_ERROR               = 1143
	//ErILLEGAL_GRANT_FOR_TABLE                 = 1144
	//ErGRANT_WRONG_HOST_OR_USER                = 1145
	//ErNO_SUCH_TABLE                           = 1146
	//ErNONEXISTING_TABLE_GRANT                 = 1147
	//ErNOT_ALLOWED_COMMAND                     = 1148
	//ErSYNTAX_ERROR                            = 1149
	//ErDELAYED_CANT_CHANGE_LOCK                = 1150
	//ErTOO_MANY_DELAYED_THREADS                = 1151
	//ErABORTING_CONNECTION                     = 1152
	//ErNET_PACKET_TOO_LARGE                    = 1153
	//ErNET_READ_ERROR_FROM_PIPE                = 1154
	//ErNET_FCNTL_ERROR                         = 1155
	//ErNET_PACKETS_OUT_OF_ORDER                = 1156
	//ErNET_UNCOMPRESS_ERROR                    = 1157
	//ErNET_READ_ERROR                          = 1158
	//ErNET_READ_INTERRUPTED                    = 1159
	//ErNET_ERROR_ON_WRITE                      = 1160
	//ErNET_WRITE_INTERRUPTED                   = 1161
	//ErTOO_LONG_STRING                         = 1162
	//ErTABLE_CANT_HANDLE_BLOB                  = 1163
	//ErTABLE_CANT_HANDLE_AUTO_INCREMENT        = 1164
	//ErDELAYED_INSERT_TABLE_LOCKED             = 1165
	//ErWRONG_COLUMN_NAME                       = 1166
	//ErWRONG_KEY_COLUMN                        = 1167
	//ErWRONG_MRG_TABLE                         = 1168
	//ErDUP_UNIQUE                              = 1169
	//ErBLOB_KEY_WITHOUT_LENGTH                 = 1170
	//ErPRIMARY_CANT_HAVE_NULL                  = 1171
	//ErTOO_MANY_ROWS                           = 1172
	//ErREQUIRES_PRIMARY_KEY                    = 1173
	//ErNO_RAID_COMPILED                        = 1174
	//ErUPDATE_WITHOUT_KEY_IN_SAFE_MODE         = 1175
	//ErKEY_DOES_NOT_EXITS                      = 1176
	//ErCHECK_NO_SUCH_TABLE                     = 1177
	//ErCHECK_NOT_IMPLEMENTED                   = 1178
	//ErCANT_DO_THIS_DURING_AN_TRANSACTION      = 1179
	//ErERROR_DURING_COMMIT                     = 1180
	//ErERROR_DURING_ROLLBACK                   = 1181
	//ErERROR_DURING_FLUSH_LOGS                 = 1182
	//ErERROR_DURING_CHECKPOINT                 = 1183
	//ErNEW_ABORTING_CONNECTION                 = 1184
	//ErDUMP_NOT_IMPLEMENTED                    = 1185
	//ErFLUSH_MASTER_BINLOG_CLOSED              = 1186
	//ErINDEX_REBUILD                           = 1187
	//ErMASTER                                  = 1188
	//ErMASTER_NET_READ                         = 1189
	//ErMASTER_NET_WRITE                        = 1190
	//ErFT_MATCHING_KEY_NOT_FOUND               = 1191
	//ErLOCK_OR_ACTIVE_TRANSACTION              = 1192
	//ErUNKNOWN_SYSTEM_VARIABLE                 = 1193
	//ErCRASHED_ON_USAGE                        = 1194
	//ErCRASHED_ON_REPAIR                       = 1195
	//ErWARNING_NOT_COMPLETE_ROLLBACK           = 1196
	//ErTRANS_CACHE_FULL                        = 1197
	//ErSLAVE_MUST_STOP                         = 1198
	//ErSLAVE_NOT_RUNNING                       = 1199
	//ErBAD_SLAVE                               = 1200
	//ErMASTER_INFO                             = 1201
	//ErSLAVE_THREAD                            = 1202
	//ErTOO_MANY_USER_CONNECTIONS               = 1203
	//ErSET_CONSTANTS_ONLY                      = 1204
	//ErLOCK_WAIT_TIMEOUT                       = 1205
	//ErLOCK_TABLE_FULL                         = 1206
	//ErREAD_ONLY_TRANSACTION                   = 1207
	//ErDROP_DB_WITH_READ_LOCK                  = 1208
	//ErCREATE_DB_WITH_READ_LOCK                = 1209
	//ErWRONG_ARGUMENTS                         = 1210
	//ErNO_PERMISSION_TO_CREATE_USER            = 1211
	//ErUNION_TABLES_IN_DIFFERENT_DIR           = 1212
	//ErLOCK_DEADLOCK                           = 1213
	//ErTABLE_CANT_HANDLE_FT                    = 1214
	//ErCANNOT_ADD_FOREIGN                      = 1215
	//ErNO_REFERENCED_ROW                       = 1216
	//ErROW_IS_REFERENCED                       = 1217
	//ErCONNECT_TO_MASTER                       = 1218
	//ErQUERY_ON_MASTER                         = 1219
	//ErERROR_WHEN_EXECUTING_COMMAND            = 1220
	//ErWRONG_USAGE                             = 1221
	//ErWRONG_NUMBER_OF_COLUMNS_IN_SELECT       = 1222
	//ErCANT_UPDATE_WITH_READLOCK               = 1223
	//ErMIXING_NOT_ALLOWED                      = 1224
	//ErDUP_ARGUMENT                            = 1225
	//ErUSER_LIMIT_REACHED                      = 1226
	//ErSPECIFIC_ACCESS_DENIED_ERROR            = 1227
	//ErLOCAL_VARIABLE                          = 1228
	//ErGLOBAL_VARIABLE                         = 1229
	//ErNO_DEFAULT                              = 1230
	//ErWRONG_VALUE_FOR_VAR                     = 1231
	//ErWRONG_TYPE_FOR_VAR                      = 1232
	//ErVAR_CANT_BE_READ                        = 1233
	//ErCANT_USE_OPTION_HERE                    = 1234
	//ErNOT_SUPPORTED_YET                       = 1235
	//ErMASTER_FATAL_ERROR_READING_BINLOG       = 1236
	//ErSLAVE_IGNORED_TABLE                     = 1237
	//ErINCORRECT_GLOBAL_LOCAL_VAR              = 1238
	//ErWRONG_FK_DEF                            = 1239
	//ErKEY_REF_DO_NOT_MATCH_TABLE_REF          = 1240
	//ErOPERAND_COLUMNS                         = 1241
	//ErSUBQUERY_NO_1_ROW                       = 1242
	//ErUNKNOWN_STMT_HANDLER                    = 1243
	//ErCORRUPT_HELP_DB                         = 1244
	//ErCYCLIC_REFERENCE                        = 1245
	//ErAUTO_CONVERT                            = 1246
	//ErILLEGAL_REFERENCE                       = 1247
	//ErDERIVED_MUST_HAVE_ALIAS                 = 1248
	//ErSELECT_REDUCED                          = 1249
	//ErTABLENAME_NOT_ALLOWED_HERE              = 1250
	//ErNOT_SUPPORTED_AUTH_MODE                 = 1251
	//ErSPATIAL_CANT_HAVE_NULL                  = 1252
	//ErCOLLATION_CHARSET_MISMATCH              = 1253
	//ErSLAVE_WAS_RUNNING                       = 1254
	//ErSLAVE_WAS_NOT_RUNNING                   = 1255
	//ErTOO_BIG_FOR_UNCOMPRESS                  = 1256
	//ErZLIB_Z_MEM_ERROR                        = 1257
	//ErZLIB_Z_BUF_ERROR                        = 1258
	//ErZLIB_Z_DATA_ERROR                       = 1259
	//ErCUT_VALUE_GROUP_CONCAT                  = 1260
	//ErWARN_TOO_FEW_RECORDS                    = 1261
	//ErWARN_TOO_MANY_RECORDS                   = 1262
	//ErWARN_NULL_TO_NOTNULL                    = 1263
	//ErWARN_DATA_OUT_OF_RANGE                  = 1264
	//WrN_DATA_TRUNCATED                        = 1265
	//ErWARN_USING_OTHER_HANDLER                = 1266
	//ErCANT_AGGREGATE_2COLLATIONS              = 1267
	//ErDROP_USER                               = 1268
	//ErREVOKE_GRANTS                           = 1269
	//ErCANT_AGGREGATE_3COLLATIONS              = 1270
	//ErCANT_AGGREGATE_NCOLLATIONS              = 1271
	//ErVARIABLE_IS_NOT_STRUCT                  = 1272
	//ErUNKNOWN_COLLATION                       = 1273
	//ErSLAVE_IGNORED_SSL_PARAMS                = 1274
	//ErSERVER_IS_IN_SECURE_AUTH_MODE           = 1275
	//ErWARN_FIELD_RESOLVED                     = 1276
	//ErBAD_SLAVE_UNTIL_COND                    = 1277
	//ErMISSING_SKIP_SLAVE                      = 1278
	//ErUNTIL_COND_IGNORED                      = 1279
	//ErWRONG_NAME_FOR_INDEX                    = 1280
	//ErWRONG_NAME_FOR_CATALOG                  = 1281
	//ErWARN_QC_RESIZE                          = 1282
	//ErBAD_FT_COLUMN                           = 1283
	//ErUNKNOWN_KEY_CACHE                       = 1284
	//ErWARN_HOSTNAME_WONT_WORK                 = 1285
	//ErUNKNOWN_STORAGE_ENGINE                  = 1286
	//ErWARN_DEPRECATED_SYNTAX                  = 1287
	//ErNON_UPDATABLE_TABLE                     = 1288
	//ErFEATURE_DISABLED                        = 1289
	//ErOPTION_PREVENTS_STATEMENT               = 1290
	//ErDUPLICATED_VALUE_IN_TYPE                = 1291
	//ErTRUNCATED_WRONG_VALUE                   = 1292
	//ErTOO_MUCH_AUTO_TIMESTAMP_COLS            = 1293
	//ErINVALID_ON_UPDATE                       = 1294
	//ErUNSUPPORTED_PS                          = 1295
	//ErGET_ERRMSG                              = 1296
	//ErGET_TEMPORARY_ERRMSG                    = 1297
	//ErUNKNOWN_TIME_ZONE                       = 1298
	//ErWARN_INVALID_TIMESTAMP                  = 1299
	//ErINVALID_CHARACTER_STRING                = 1300
	//ErWARN_ALLOWED_PACKET_OVERFLOWED          = 1301
	//ErCONFLICTING_DECLARATIONS                = 1302
	//ErSP_NO_RECURSIVE_CREATE                  = 1303
	//ErSP_ALREADY_EXISTS                       = 1304
	//ErSP_DOES_NOT_EXIST                       = 1305
	//ErSP_DROP_FAILED                          = 1306
	//ErSP_STORE_FAILED                         = 1307
	//ErSP_LILABEL_MISMATCH                     = 1308
	//ErSP_LABEL_REDEFINE                       = 1309
	//ErSP_LABEL_MISMATCH                       = 1310
	//ErSP_UNINIT_VAR                           = 1311
	//ErSP_BADSELECT                            = 1312
	//ErSP_BADRETURN                            = 1313
	//ErSP_BADSTATEMENT                         = 1314
	//ErUPDATE_LOG_DEPRECATED_IGNORED           = 1315
	//ErUPDATE_LOG_DEPRECATED_TRANSLATED        = 1316
	//ErQUERY_INTERRUPTED                       = 1317
	//ErSP_WRONG_NO_OF_ARGS                     = 1318
	//ErSP_COND_MISMATCH                        = 1319
	//ErSP_NORETURN                             = 1320
	//ErSP_NORETURNEND                          = 1321
	//ErSP_BAD_CURSOR_QUERY                     = 1322
	//ErSP_BAD_CURSOR_SELECT                    = 1323
	//ErSP_CURSOR_MISMATCH                      = 1324
	//ErSP_CURSOR_ALREADY_OPEN                  = 1325
	//ErSP_CURSOR_NOT_OPEN                      = 1326
	//ErSP_UNDECLARED_VAR                       = 1327
	//ErSP_WRONG_NO_OF_FETCH_ARGS               = 1328
	//ErSP_FETCH_NO_DATA                        = 1329
	//ErSP_DUP_PARAM                            = 1330
	//ErSP_DUP_VAR                              = 1331
	//ErSP_DUP_COND                             = 1332
	//ErSP_DUP_CURS                             = 1333
	//ErSP_CANT_ALTER                           = 1334
	//ErSP_SUBSELECT_NYI                        = 1335
	//ErSTMT_NOT_ALLOWED_IN_SF_OR_TRG           = 1336
	//ErSP_VARCOND_AFTER_CURSHNDLR              = 1337
	//ErSP_CURSOR_AFTER_HANDLER                 = 1338
	//ErSP_CASE_NOT_FOUND                       = 1339
	//ErFPARSER_TOO_BIG_FILE                    = 1340
	//ErFPARSER_BAD_HEADER                      = 1341
	//ErFPARSER_EOF_IN_COMMENT                  = 1342
	//ErFPARSER_ERROR_IN_PARAMETER              = 1343
	//ErFPARSER_EOF_IN_UNKNOWN_PARAMETER        = 1344
	//ErVIEW_NO_EXPLAIN                         = 1345
	//ErFRM_UNKNOWN_TYPE                        = 1346
	//ErWRONG_OBJECT                            = 1347
	//ErNONUPDATEABLE_COLUMN                    = 1348
	//ErVIEW_SELECT_DERIVED                     = 1349
	//ErVIEW_SELECT_CLAUSE                      = 1350
	//ErVIEW_SELECT_VARIABLE                    = 1351
	//ErVIEW_SELECT_TMPTABLE                    = 1352
	//ErVIEW_WRONG_LIST                         = 1353
	//ErWARN_VIEW_MERGE                         = 1354
	//ErWARN_VIEW_WITHOUT_KEY                   = 1355
	//ErVIEW_INVALID                            = 1356
	//ErSP_NO_DROP_SP                           = 1357
	//ErSP_GOTO_IN_HNDLR                        = 1358
	//ErTRG_ALREADY_EXISTS                      = 1359
	//ErTRG_DOES_NOT_EXIST                      = 1360
	//ErTRG_ON_VIEW_OR_TEMP_TABLE               = 1361
	//ErTRG_CANT_CHANGE_ROW                     = 1362
	//ErTRG_NO_SUCH_ROW_IN_TRG                  = 1363
	//ErNO_DEFAULT_FOR_FIELD                    = 1364
	//ErDIVISION_BY_ZERO                        = 1365
	//ErTRUNCATED_WRONG_VALUE_FOR_FIELD         = 1366
	//ErILLEGAL_VALUE_FOR_TYPE                  = 1367
	//ErVIEW_NONUPD_CHECK                       = 1368
	//ErVIEW_CHECK_FAILED                       = 1369
	//ErPROCACCESS_DENIED_ERROR                 = 1370
	//ErRELAY_LOG_FAIL                          = 1371
	//ErPASSWD_LENGTH                           = 1372
	//ErUNKNOWN_TARGET_BINLOG                   = 1373
	//ErIO_ERR_LOG_INDEX_READ                   = 1374
	//ErBINLOG_PURGE_PROHIBITED                 = 1375
	//ErFSEEK_FAIL                              = 1376
	//ErBINLOG_PURGE_FATAL_ERR                  = 1377
	//ErLOG_IN_USE                              = 1378
	//ErLOG_PURGE_UNKNOWN_ERR                   = 1379
	//ErRELAY_LOG_INIT                          = 1380
	//ErNO_BINARY_LOGGING                       = 1381
	//ErRESERVED_SYNTAX                         = 1382
	//ErWSAS_FAILED                             = 1383
	//ErDIFF_GROUPS_PROC                        = 1384
	//ErNO_GROUP_FOR_PROC                       = 1385
	//ErORDER_WITH_PROC                         = 1386
	//ErLOGGING_PROHIBIT_CHANGING_OF            = 1387
	//ErNO_FILE_MAPPING                         = 1388
	//ErWRONG_MAGIC                             = 1389
	//ErPS_MANY_PARAM                           = 1390
	//ErKEY_PART_0                              = 1391
	//ErVIEW_CHECKSUM                           = 1392
	//ErVIEW_MULTIUPDATE                        = 1393
	//ErVIEW_NO_INSERT_FIELD_LIST               = 1394
	//ErVIEW_DELETE_MERGE_VIEW                  = 1395
	//ErCANNOT_USER                             = 1396
	//ErXAER_NOTA                               = 1397
	//ErXAER_INVAL                              = 1398
	//ErXAER_RMFAIL                             = 1399
	//ErXAER_OUTSIDE                            = 1400
	//ErXAER_RMERR                              = 1401
	//ErXA_RBROLLBACK                           = 1402
	//ErNONEXISTING_PROC_GRANT                  = 1403
	//ErPROC_AUTO_GRANT_FAIL                    = 1404
	//ErPROC_AUTO_REVOKE_FAIL                   = 1405
	//ErDATA_TOO_LONG                           = 1406
	//ErSP_BAD_SQLSTATE                         = 1407
	//ErSTARTUP                                 = 1408
	//ErLOAD_FROM_FIXED_SIZE_ROWS_TO_VAR        = 1409
	//ErCANT_CREATE_USER_WITH_GRANT             = 1410
	//ErWRONG_VALUE_FOR_TYPE                    = 1411
	//ErTABLE_DEF_CHANGED                       = 1412
	//ErSP_DUP_HANDLER                          = 1413
	//ErSP_NOT_VAR_ARG                          = 1414
	//ErSP_NO_RETSET                            = 1415
	//ErCANT_CREATE_GEOMETRY_OBJECT             = 1416
	//ErFAILED_ROUTINE_BREAK_BINLOG             = 1417
	//ErBINLOG_UNSAFE_ROUTINE                   = 1418
	//ErBINLOG_CREATE_ROUTINE_NEED_SUPER        = 1419
	//ErEXEC_STMT_WITH_OPEN_CURSOR              = 1420
	//ErSTMT_HAS_NO_OPEN_CURSOR                 = 1421
	//ErCOMMIT_NOT_ALLOWED_IN_SF_OR_TRG         = 1422
	//ErNO_DEFAULT_FOR_VIEW_FIELD               = 1423
	//ErSP_NO_RECURSION                         = 1424
	//ErTOO_BIG_SCALE                           = 1425
	//ErTOO_BIG_PRECISION                       = 1426
	//ErM_BIGGER_THAN_D                         = 1427
	//ErWRONG_LOCK_OF_SYSTEM_TABLE              = 1428
	//ErCONNECT_TO_FOREIGN_DATA_SOURCE          = 1429
	//ErQUERY_ON_FOREIGN_DATA_SOURCE            = 1430
	//ErFOREIGN_DATA_SOURCE_DOESNT_EXIST        = 1431
	//ErFOREIGN_DATA_STRING_INVALID_CANT_CREATE = 1432
	//ErFOREIGN_DATA_STRING_INVALID             = 1433
	//ErCANT_CREATE_FEDERATED_TABLE             = 1434
	//ErTRG_IN_WRONG_SCHEMA                     = 1435
	//ErSTACK_OVERRUN_NEED_MORE                 = 1436
	//ErTOO_LONG_BODY                           = 1437
	//ErWARN_CANT_DROP_DEFAULT_KEYCACHE         = 1438
	//ErTOO_BIG_DISPLAYWIDTH                    = 1439
	//ErXAER_DUPID                              = 1440
	//ErDATETIME_FUNCTION_OVERFLOW              = 1441
	//ErCANT_UPDATE_USED_TABLE_IN_SF_OR_TRG     = 1442
	//ErVIEW_PREVENT_UPDATE                     = 1443
	//ErPS_NO_RECURSION                         = 1444
	//ErSP_CANT_SET_AUTOCOMMIT                  = 1445
	//ErMALFORMED_DEFINER                       = 1446
	//ErVIEW_FRM_NO_USER                        = 1447
	//ErVIEW_OTHER_USER                         = 1448
	//ErNO_SUCH_USER                            = 1449
	//ErFORBID_SCHEMA_CHANGE                    = 1450
	//ErROW_IS_REFERENCED_2                     = 1451
	//ErNO_REFERENCED_ROW_2                     = 1452
	//ErSP_BAD_VAR_SHADOW                       = 1453
	//ErTRG_NO_DEFINER                          = 1454
	//ErOLD_FILE_FORMAT                         = 1455
	//ErSP_RECURSION_LIMIT                      = 1456
	//ErSP_PROC_TABLE_CORRUPT                   = 1457
	//ErSP_WRONG_NAME                           = 1458
	//ErTABLE_NEEDS_UPGRADE                     = 1459
	//ErSP_NO_AGGREGATE                         = 1460
	//ErMAX_PREPARED_STMT_COUNT_REACHED         = 1461
	//ErVIEW_RECURSIVE                          = 1462
	//ErNON_GROUPING_FIELD_USED                 = 1463
	//ErTABLE_CANT_HANDLE_SPKEYS                = 1464
	//ErNO_TRIGGERS_ON_SYSTEM_SCHEMA            = 1465
	//ErREMOVED_SPACES                          = 1466
	//ErAUTOINC_READ_FAILED                     = 1467
	//ErUSERNAME                                = 1468
	//ErHOSTNAME                                = 1469
	//ErWRONG_STRING_LENGTH                     = 1470
	//ErNON_INSERTABLE_TABLE                    = 1471
)

