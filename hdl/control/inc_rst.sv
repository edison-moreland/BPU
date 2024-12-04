module inc_rst #(
    parameter N = 3
) (
    input  logic         rst,
    input  logic [N-1:0] din,
    output logic [N-1:0] dout,
);
    logic [N-1:0] dinInc;
    logic overflowInc;

    always_comb begin
        if (rst) dout = 0;
        else dout = din + 1;
    end
endmodule