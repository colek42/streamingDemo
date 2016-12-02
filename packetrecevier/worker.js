importScripts('ffmpeg-h264.js');


i = 0;

initDecoder();

onmessage = function(e) {
    console.log(e.data);
    if (e.data.topic) {

        switch (e.data.topic) {

            case "openUri":
                console.log("Starting Decode on Uri: " + e.data.topic);
                startDecode(e.data.data);
                break;
        }
    }
};

function startDecode(uri) {
    initWebsocket(uri);
}

function initDecoder() {
    Module.ccall('avcodec_register_all');
    // find h264 decoder
    codec = Module.ccall('avcodec_find_decoder_by_name', 'number', ['string'], ["h264"]);
    if (codec === 0)
        alert("Could not find H264 codec ");

    ctx = Module.ccall('avcodec_alloc_context3', 'number', ['number'], [codec]);
    ret = Module.ccall('avcodec_open2', 'number', ['number', 'number', 'number'], [ctx, codec, 0]);
    if (ret < 0)
        alert("Could not open codec ");

    // allocate packet
    pkt = Module._malloc(96);
    Module.ccall('av_init_packet', 'null', ['number'], [pkt]);
    pktData = Module._malloc(1024 * 3000);
    Module.setValue(pkt + 24, pktData, '*');
    // allocate video frame
    frame = Module.ccall('avcodec_alloc_frame', 'number');
    if (!frame)
        alert("Could not allocate video frame ");

    // init decode frame function
    new_packet = Module.cwrap('av_packet_from_data', 'number', ['number', 'number', 'number']);
    decode_frame = Module.cwrap('avcodec_decode_video2', 'number', ['number', 'number', 'number', 'number']);
    got_frame = Module._malloc(4);
}

function initWebsocket(uri) {
    ws = new WebSocket(uri);
    ws.binaryType = 'arraybuffer';
    console.log("Opening ws  " + uri);
    ws.onmessage = decodePkt;
}

function decodePkt(event) {
    var decodedFrame = decode(event.data);
    if (decodedFrame) {
        self.postMessage(decodedFrame, [
            decodedFrame.frameYData.buffer,
            decodedFrame.frameUData.buffer,
            decodedFrame.frameVData.buffer
        ]);
    }
}

function decode(data) {
    var buffer = new Uint8Array(data);
    var dataSize = data.byteLength;
    Module.setValue(pkt + 28, dataSize, 'i32');
    Module.writeArrayToMemory(buffer, pktData);

    var len = decode_frame(ctx, frame, got_frame, pkt);

    if (len < 0) {
        console.log("Error while decoding frame");
        return;
    }

    if (Module.getValue(got_frame, 'i8') === 0) {
        console.log("No frame");
        return;
    }

    var decoded_frame = frame;
    var frame_width = Module.getValue(decoded_frame + 68, 'i32');
    var frame_height = Module.getValue(decoded_frame + 72, 'i32');

    // copy Y channel to canvas
    var frameYDataPtr = Module.getValue(decoded_frame, '*');
    var frameUDataPtr = Module.getValue(decoded_frame + 4, '*');
    var frameVDataPtr = Module.getValue(decoded_frame + 8, '*');


    return {
        frame_width: frame_width,
        frame_height: frame_height,
        frameYDataPtr: frameYDataPtr,
        frameUDataPtr: frameUDataPtr,
        frameVDataPtr: frameVDataPtr,
        frameYData: new Uint8Array(Module.HEAPU8.buffer.slice(frameYDataPtr, frameYDataPtr + frame_width * frame_height)),
        frameUData: new Uint8Array(Module.HEAPU8.buffer.slice(frameUDataPtr, frameUDataPtr + frame_width / 2 * frame_height / 2)),
        frameVData: new Uint8Array(Module.HEAPU8.buffer.slice(frameVDataPtr, frameVDataPtr + frame_width / 2 * frame_height / 2))
    };
}
