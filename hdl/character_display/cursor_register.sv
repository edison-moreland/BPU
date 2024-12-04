`include "inc_dec_rst.sv"

module cursor_register #(
    parameter N = 3
) (
    input  logic         clk, inc, dec, rst,
    input  logic [N-1:0] din,
    output logic [N-1:0] dout,
);

    logic [N-1:0] did;
    inc_dec_rst id(
        .inc(inc), .dec(dec), .rst(rst),
        .din(din), .dout(did)
    );

    always_ff @(negedge clk) begin
        dout <= did;
    end

endmodule