module incrementer #(
    parameter N = 3
) (
    // input  logic increment_in,
    input  logic [N-1:0] data_in,
    // output logic increment_out,
    output logic [N-1:0] data_out,
);

    assign data_out = data_in + 1;

endmodule